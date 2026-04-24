package middleware

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"golang.org/x/time/rate"
)

// 限流器, 每秒 5 次，桶容量 10 (平均每秒允许 5 次请求，最多可瞬间突发 10 次请求)
func RateLimit() middleware.Middleware {
	
	//登录接口 rate.NewLimiter(1, 3)
	//后台API rate.NewLimiter(5, 10)
	//搜索接口 rate.NewLimiter(10, 20)
	//上传接口 rate.NewLimiter(1, 2)
	//内部服务调用（微服务）rate.NewLimiter(50, 100)
	limiter := rate.NewLimiter(5, 10)

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