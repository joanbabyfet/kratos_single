package middleware

import (
	"context"
	"kratos_single/internal/pkg/auth"
	"strings"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"

	kratosmd "github.com/go-kratos/kratos/v2/metadata"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

// AuthMiddleware 认证中间件, 请求参数没token则通过, 有则解析
func AuthMiddleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			
			var token string

			// HTTP Header 获取 JWT, Authorization: Bearer xxxxxx
			if r, ok := khttp.RequestFromServerContext(ctx); ok {
				token = r.Header.Get("Authorization")
			}
			
			//gRPC / 内部服务 metadata 获取 JWT
			if token == "" {
				if md, ok := kratosmd.FromServerContext(ctx); ok {
					token = md.Get("authorization")
				}
			}

			// 清理 Bearer
			token = strings.TrimSpace(token)
			token = strings.TrimPrefix(token, "Bearer ")
			token = strings.TrimSpace(token)
			
			// 没 token 允许匿名访问
			if token == "" {
				return handler(ctx, req)
			}

			// 解析 JWT
			payload, err := auth.ValidateToken(token)
			if err != nil {
				return nil, errors.Unauthorized(
					"UNAUTHORIZED",
					"登录已过期",
				)
			}
			
			// 写入 ctx（后续 service/biz 可直接取）
			ctx = auth.SetUser(ctx, payload.UserID)
			ctx = auth.SetRole(ctx, payload.Role)
			
			//服务间继续透传（调用下游服务）
			ctx = kratosmd.AppendToClientContext(
				ctx,
				"authorization", "Bearer "+token,
				"x-user-id", payload.UserID,
				"x-role", payload.Role,
			)

			//测试用
			// md, ok := kratosmd.FromClientContext(ctx)
			// if ok {
			// 	log.Println("client metadata:", md)
			// }

			return handler(ctx, req)
		}
	}
}