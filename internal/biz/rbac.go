package biz

import (
	"context"
)

//对外业务接口（给 middleware / service 用）
type RBAC interface {
	HasPermission(ctx context.Context, uid string, perm string) (bool, error)
}

//对内repo 接口（biz 依赖 data）
type RBACRepo interface {
	HasPermission(ctx context.Context, uid string, perm string) (bool, error)
}

type RBACUsecase struct {
	repo RBACRepo
}

//构造函数
func NewRBACUsecase(repo RBACRepo) *RBACUsecase {
    return &RBACUsecase{repo: repo}
}

func (r *RBACUsecase) HasPermission(ctx context.Context, uid string, perm string) (bool, error) {
	return r.repo.HasPermission(ctx, uid, perm)
}