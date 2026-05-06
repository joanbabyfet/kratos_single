package data

import (
	"context"
	"fmt"
	"kratos_single/internal/biz"
	"time"

	"github.com/redis/go-redis/v9"
)

type rbacRepo struct {
	data *Data
	rds *redis.Client
}

//构造函数
func NewRBACRepo(data *Data, rds *redis.Client) biz.RBACRepo {
	return &rbacRepo{
		data: data,
		rds: rds,
	}
}

//实现接口
func (r *rbacRepo) HasPermission(ctx context.Context, uid string, perm string) (bool, error) {

	//查 Redis
	val, err := r.rds.Get(ctx, r.cacheKey(uid, perm)).Result()
	if err == nil {
		return val == "1", nil
	}

	// Redis 真异常（不是 miss）
	if err != redis.Nil {
		return false, err
	}

	//查 MySQL
	var count int64

	err = r.data.db.Raw(`
		SELECT COUNT(1)
		FROM kk_admin a
		JOIN kk_role_permission rp ON a.role_id = rp.role_id
		JOIN kk_permission p ON rp.permission_id = p.id
		WHERE a.id = ?
		AND p.perms = ?
	`, uid, perm).Scan(&count).Error
	if err != nil {
		return false, err
	}

	ok := count > 0

	//回写 Redis
	cacheVal := "0"
	if ok {
		cacheVal = "1"
	}

	_ = r.rds.Set(ctx, r.cacheKey(uid, perm), cacheVal, time.Hour).Err()

	return ok, nil
}

//获取缓存键
func (r *rbacRepo) cacheKey(uid string, perm string) string {
    return fmt.Sprintf("rbac:%s:%s", uid, perm)
}