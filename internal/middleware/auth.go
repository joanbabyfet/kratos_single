package middleware

import (
	"context"
	"strings"

	"github.com/go-kratos/kratos/v2/errors"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport"

	"kratos_single/internal/pkg/auth"

	kratosmd "github.com/go-kratos/kratos/v2/metadata"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

//
// =========================
// AuthMiddleware
// 非白名单接口：验证登录并写入 ctx
// =========================
//

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
			// if token == "" {
			// 	return handler(ctx, req)
			// }
			if token == "" {
				return nil, errors.Unauthorized(
					"UNAUTHORIZED",
					"请先登录",
				)
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

//
// =========================
// AdminAuthMiddleware
// 后台接口：admin 权限校验
// =========================
//

func AdminAuthMiddleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {

			uid := auth.GetUser(ctx)
			if uid == "" {
				return nil, kerrors.Unauthorized(
					"UNAUTHORIZED",
					"未登录",
				)
			}

			role := auth.GetRole(ctx)
			if role != "admin" {
				return nil, kerrors.Forbidden(
					"FORBIDDEN",
					"无后台权限",
				)
			}

			return handler(ctx, req)
		}
	}
}

//
// =========================
// ClientAuthMiddleware
// 前台接口：普通登录校验
// =========================
//

func ClientAuthMiddleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {

			uid := auth.GetUser(ctx)
			if uid == "" {
				return nil, kerrors.Unauthorized(
					"UNAUTHORIZED",
					"未登录",
				)
			}

			return handler(ctx, req)
		}
	}
}

//
// =========================
// 白名单（跳过登录）
// =========================
//
var whiteList = map[string]struct{}{
	"/api.admin.v1.Admin/Login":	{},
	"/api.client.v1.User/Login":	{},
	"/api.client.v1.User/Register":	{},
}
// 返回 true = 走 Auth
// 返回 false = 不走 Auth
func WhiteListMatcher() func(ctx context.Context, operation string) bool {
	return func(ctx context.Context, operation string) bool {
		
		// operation 类似: /api.client.v1.User/Login
		// HTTP / gRPC 场景统一从 transport 取真实路由

		if tr, ok := transport.FromServerContext(ctx); ok {
			path := tr.Operation()
			for k := range whiteList {
				if strings.Contains(path, k) {
					return false
				}
			}
		}

		return true
	}
}

//
// =========================
// admin 接口匹配
// =========================
//

func AdminMatcher() selector.MatchFunc {
	return func(ctx context.Context, operation string) bool {
		// 白名单接口不走 admin 登录校验
		if _, ok := whiteList[operation]; ok {
			return false
		}

		return strings.HasPrefix(operation, "/api.admin.v1.")
	}
}

//
// =========================
// client 接口匹配
// =========================
//

func ClientMatcher() selector.MatchFunc {
	return func(ctx context.Context, operation string) bool {
		// 白名单接口不走 client 登录校验
		if _, ok := whiteList[operation]; ok {
			return false
		}

		return strings.HasPrefix(operation, "/api.client.v1.")
	}
}