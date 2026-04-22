package middleware

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"golang.org/x/time/rate"
)

func RateLimit() middleware.Middleware {
	// 每秒 10 次，桶容量 20 (平均每秒允许 10 次请求，最多可瞬间突发 20 次请求)
	limiter := rate.NewLimiter(10, 20)

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {

			// 超过限制
			if !limiter.Allow() {
				return nil, errors.New(
					429,
					"RATE_LIMIT",
					"请求过于频繁",
				)
			}

			return handler(ctx, req)
		}
	}
}