package data

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"kratos_single/internal/biz"
	"kratos_single/internal/data/model"
	"kratos_single/internal/pkg/utils"

	"math/rand"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/jinzhu/copier"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type adminRepo struct {
	data *Data
	rdb *redis.Client
	log *log.Helper
}

//构造函数注入
func NewAdminRepo(data *Data, rdb *redis.Client, logger log.Logger) biz.AdminRepo {
	return &adminRepo{
		data: data,
		rdb: rdb,
		log: log.NewHelper(logger),
	}
}

// 获取列表
func (r *adminRepo) List(ctx context.Context, q *biz.AdminQuery) ([]*biz.Admin, int64, error) {

	var list []model.Admin
	var total int64

	//绑定请求上下文, 指定操作表
	db := r.data.db.WithContext(ctx).
		Model(&model.Admin{}).
		Where("delete_time = 0")

	if q.Username != "" {
		db = db.Where("title LIKE ?", "%"+q.Username+"%")
	}
	if q.Realname != "" {
		db = db.Where("title LIKE ?", "%"+q.Realname+"%")
	}
	if q.Email != "" {
		db = db.Where("title LIKE ?", "%"+q.Email+"%")
	}

	if q.Limit > 0 {
		//返回前几条
		err := db.
			Order("create_time DESC").
			Limit(q.Limit).
			Find(&list).Error

		if err != nil {
			return nil, 0, err
		}

	} else {
		//分页模式
		db.Count(&total)

		err := db.
			Order("create_time DESC").
			Offset((q.Page - 1) * q.PageSize).
			Limit(q.PageSize).
			Find(&list).Error

		if err != nil {
			return nil, 0, err
		}
	}

	var res []*biz.Admin
	for _, v := range list {
		var item biz.Admin
		_ = copier.Copy(&item, &v) // 把 model 里的字段复制到 biz, 新增字段不用改代码
		res = append(res, &item)
	}

	return res, total, nil
}

// 获取详情
func (r *adminRepo) GetById(ctx context.Context, id string, useCache bool) (*biz.Admin, error) {
	
	key := KeyUserDetail(id)

	// 只有允许缓存才读 Redis
	if useCache {
		val, err := r.rdb.Get(ctx, key).Result()
		if err == nil {
			//防穿透时写入null, 这里不会再查库直接返回
			if val == "null" {
				return nil, gorm.ErrRecordNotFound
			}
			var res biz.Admin
			if err := json.Unmarshal([]byte(val), &res); err == nil {
				return &res, nil
			}
		}
	}

	//查数据库
	var m model.Admin
	err := r.data.db.WithContext(ctx).
		Where("id = ? AND delete_time = 0", id).
		First(&m).Error

	if err != nil {
		// 防穿透
		if useCache && err == gorm.ErrRecordNotFound {
			_ = r.rdb.Set(ctx, key, "null", time.Minute).Err()
		}
		return nil, err
	}

	var res biz.Admin
	// 把 model 里的字段复制到 biz, 新增字段不用改代码
	if err := copier.Copy(&res, &m); err != nil {
		return nil, err
	}

	// 只有允许缓存才写 Redis
	if useCache {
		bytes, err := json.Marshal(res)
		if err == nil {
			//_ = r.rdb.Set(ctx, key, bytes, 10*time.Minute).Err()
			//加随机过期（防雪崩）
			expire := 10*time.Minute + time.Duration(rand.Intn(300))*time.Second
			_ = r.rdb.Set(ctx, key, bytes, expire).Err()
		}
	}

	return &res, nil
}

func (r *adminRepo) GetByUsername(ctx context.Context, username string) (*biz.Admin, error) {

	var m model.Admin

	err := r.data.db.WithContext(ctx).
		Model(&model.Admin{}).
		Where("username = ? AND delete_time = 0", username).
		First(&m).Error

	if err != nil {
		return nil, err
	}

	var u biz.Admin

	// model -> biz
	if err := copier.Copy(&u, &m); err != nil {
		return nil, err
	}

	return &u, nil
}

// 添加
func (r *adminRepo) Create(ctx context.Context, a *biz.Admin) (string, error) {

	var id string

	now := utils.Timestamp()

	err := r.data.WithTx(ctx, func(tx *gorm.DB) error {

		var m model.Admin

		// copy biz → model
		if err := copier.Copy(&m, a); err != nil {
			return err
		}

		// 补充 DB 字段
		m.Status = 1
		m.CreateTime = now

		if err := tx.Create(&m).Error; err != nil {
			return err
		}

		id = m.Id
		return nil
	})

	if err != nil {
		return "", err
	}

	return id, nil
}

// 修改, 一律删除缓存（统一处理）
func (r *adminRepo) Update(ctx context.Context, a *biz.Admin, lockMode string) error {
	key := KeyUserDetail(a.Id)

	update := map[string]interface{}{
		"realname":    a.Realname,
		"sex":         a.Sex,
		"email":       a.Email,
		"update_time": utils.Timestamp(),
		"update_user": a.UpdateUser,
	}

	//先更新数据库（事务）
	err := r.data.WithTx(ctx, func(tx *gorm.DB) error {
		// ===============================
		// 无锁更新
		// ===============================
		if lockMode == "" {
			return tx.Model(&model.Admin{}).
				Where("id = ? AND delete_time = 0", a.Id).
				Updates(update).Error
		}

		var admin model.Admin
		db := tx

		// ===============================
		// FOR UPDATE
		// ===============================
		if lockMode == "update" {
			db = db.Clauses(clause.Locking{
				Strength: "UPDATE",
			})
		}

		// ===============================
		// SHARE MODE
		// MySQL8 = FOR SHARE
		// ===============================
		if lockMode == "share" {
			db = db.Clauses(clause.Locking{
				Strength: "SHARE",
			})
		}

		// 先锁记录
		if err := db.
			Where("id = ? AND delete_time = 0", a.Id).
			First(&admin).Error; err != nil {
			return err
		}

		// 再更新
		return tx.Model(&admin).
			Updates(update).Error
	})

	// 成功后删除缓存
	if err == nil {
		if err := r.rdb.Del(ctx, key).Err(); err != nil {
			//写入日志
			r.log.Errorw(
				"msg", "cache update failed",
				"id", a.Id,
				"err", err,
			)
		}
	}

	return err
}

//修改密码
func (r *adminRepo) UpdatePassword(ctx context.Context, id string, password string) error {

	update := map[string]interface{}{
		"password":    password,
		"update_time": utils.Timestamp(),
		"update_user": "",
	}

	return r.data.db.WithContext(ctx).
		Model(&model.Admin{}).
		Where("id = ? AND delete_time = 0", id).
		Updates(update).Error
}

//更新登录信息
func (r *adminRepo) UpdateLogin(ctx context.Context, id string, ip string) error {

	update := map[string]interface{}{
		"login_time":  utils.Timestamp(),
		"login_ip":    ip,
		"update_time": utils.Timestamp(),
		"update_user": "",
	}

	return r.data.db.WithContext(ctx).
		Model(&model.Admin{}).
		Where("id = ? AND delete_time = 0", id).
		Updates(update).Error
}

// Delete（软删除）, 一律删除缓存（统一处理）
func (r *adminRepo) Delete(ctx context.Context, id string, user string) error {

	key := KeyUserDetail(id)

	//事务
	err := r.data.WithTx(ctx, func(tx *gorm.DB) error {
		return tx.Model(&model.Admin{}).
			Where("id = ? AND delete_time = 0", id).
			Updates(map[string]interface{}{
				"delete_time": time.Now().Unix(),
				"delete_user": user,
			}).Error
	})

	// 事务成功后，再删缓存
	if err == nil {
		if err := r.rdb.Del(ctx, key).Err(); err != nil {
			//写入日志
			r.log.Errorw(
				"msg", "cache del failed",
				"id", id,
				"err", err,
			)
		}
	}

	return err
}

//获取详情缓存键
func KeyAdminDetail(id string) string {
	return fmt.Sprintf("admin:%s", id)
}