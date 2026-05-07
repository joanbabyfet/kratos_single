package middleware

import (
	"context"

	"kratos_single/internal/biz"
	"kratos_single/internal/pkg/auth"

	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

//权限映射表(写死) 这里为准再同步到表
var permMap = map[string]string{
	"/api.admin.v1.AdminArticle/CreateArticle": "article:create",
	"/api.admin.v1.AdminArticle/UpdateArticle": "article:update",
}

func RBACMiddleware(rbac biz.RBAC) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {

			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return handler(ctx, req)
			}
			
			//路径示例 /api.admin.v1.AdminArticle/CreateArticle
			path := tr.Operation()
			
			//获取权限标识
			perm, ok := permMap[path]
			if !ok {
				return handler(ctx, req)
			}

			uid := auth.GetUser(ctx)
			if uid == "" {
				return nil, kerrors.Unauthorized("UNAUTHORIZED", "未登录")
			}

			pass, err := rbac.HasPermission(ctx, uid, perm)
			if err != nil {
				return nil, kerrors.InternalServer("RBAC_ERROR", err.Error())
			}
			if !pass {
				return nil, kerrors.Forbidden("FORBIDDEN", "无权限")
			}

			return handler(ctx, req)
		}
	}
}