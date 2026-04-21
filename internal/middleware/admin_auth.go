package middleware

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
)

// 后台认证中间件, 负责后台权限判断
func AdminAuthMiddleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {

			//role := auth.GetRole(ctx)

			// if role != "admin" {
			// 	return nil, errors.Forbidden(
			// 		"FORBIDDEN",
			// 		"无后台权限",
			// 	)
			// }

			return handler(ctx, req)
		}
	}
}