package middleware

import (
	"context"
	"strings"

	"github.com/go-kratos/kratos/v2/transport"
)

//按 operation 匹配，不要按 /v1/api/login HTTP path
var whiteList = map[string]struct{}{
	"/api.client.v1.User/Login":    {},
	"/api.client.v1.User/Register": {},
}

// 返回 true = 走 Auth
// 返回 false = 不走 Auth
func WhiteListMatcher() func(ctx context.Context, operation string) bool {
	return func(ctx context.Context, operation string) bool {
		
		// operation 类似: /api.user.v1.User/Login
		// HTTP path 可取 transport

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