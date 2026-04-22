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

type AdRepo struct {
	data *Data
	rdb *redis.Client
	log *log.Helper
}

//构造函数注入
func NewAdRepo(data *Data, rdb *redis.Client, logger log.Logger) biz.AdRepo {
	return &AdRepo{
		data: data,
		rdb: rdb,
		log: log.NewHelper(logger),
	}
}

// 获取列表
func (r *AdRepo) List(ctx context.Context, q *biz.AdQuery) ([]*biz.Ad, int64, error) {

	var list []model.Ad
	var total int64

	//绑定请求上下文, 指定操作表
	db := r.data.db.WithContext(ctx).
		Model(&model.Ad{}).
		Where("delete_time = 0")

	if q.Title != "" {
		db = db.Where("title LIKE ?", "%"+q.Title+"%")
	}
	if q.Catid > 0 {
		db = db.Where("catid = ?", q.Catid)
	}

	if q.Limit > 0 {
		//返回前几条
		err := db.
			Order("sort ASC, id DESC").
			Limit(q.Limit).
			Find(&list).Error

		if err != nil {
			return nil, 0, err
		}

	} else {
		//分页模式
		db.Count(&total)

		err := db.
			Order("sort ASC, id DESC").
			Offset((q.Page - 1) * q.PageSize).
			Limit(q.PageSize).
			Find(&list).Error

		if err != nil {
			return nil, 0, err
		}
	}

	var res []*biz.Ad
	for _, v := range list {
		var item biz.Ad
		_ = copier.Copy(&item, &v) // 把 model 里的字段复制到 biz, 新增字段不用改代码
		res = append(res, &item)
	}

	return res, total, nil
}

// 获取详情
func (r *AdRepo) GetById(ctx context.Context, id int64, useCache bool) (*biz.Ad, error) {
	
	key := KeyAdDetail(id)

	// 只有允许缓存才读 Redis
	if useCache {
		val, err := r.rdb.Get(ctx, key).Result()
		if err == nil {
			//防穿透时写入null, 这里不会再查库直接返回
			if val == "null" {
				return nil, gorm.ErrRecordNotFound
			}
			var res biz.Ad
			if err := json.Unmarshal([]byte(val), &res); err == nil {
				return &res, nil
			}
		}
	}

	//查数据库
	var m model.Ad
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

	var res biz.Ad
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

// 添加
func (r *AdRepo) Create(ctx context.Context, a *biz.Ad) (int64, error) {

	var id int64

	now := utils.Timestamp()

	err := r.data.WithTx(ctx, func(tx *gorm.DB) error {

		var m model.Ad

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

		id = int64(m.Id)
		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}

// 修改, 一律删除缓存（统一处理）
func (r *AdRepo) Update(ctx context.Context, a *biz.Ad, lockMode string) error {
	key := KeyAdDetail(int64(a.Id))

	update := map[string]interface{}{
		"catid":       a.Catid,
		"title":       a.Title,
		"img":         a.Img,
		"url":     	   a.Url,
		"sort":        a.Sort,
		"status":      a.Status,
		"update_time": utils.Timestamp(),
		"update_user": a.UpdateUser,
	}

	//先更新数据库（事务）
	err := r.data.WithTx(ctx, func(tx *gorm.DB) error {
		// ===============================
		// 无锁更新
		// ===============================
		if lockMode == "" {
			return tx.Model(&model.Ad{}).
				Where("id = ? AND delete_time = 0", a.Id).
				Updates(update).Error
		}

		var ad model.Ad
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
			First(&ad).Error; err != nil {
			return err
		}

		// 再更新
		return tx.Model(&ad).
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

// Delete（软删除）, 一律删除缓存（统一处理）
func (r *AdRepo) Delete(ctx context.Context, id int64, user string) error {

	key := KeyAdDetail(id)

	//事务
	err := r.data.WithTx(ctx, func(tx *gorm.DB) error {
		return tx.Model(&model.Ad{}).
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
func KeyAdDetail(id int64) string {
	return fmt.Sprintf("ad:%d", id)
}