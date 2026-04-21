package middleware

import (
	"context"
	"strings"
)

// 返回 true = 走 Auth
// 返回 false = 不走 Auth
func AdminMatcher() func(ctx context.Context, operation string) bool {
	return func(ctx context.Context, operation string) bool {
		return strings.HasPrefix(operation, "/api.admin.v1.")
	}
}