package middleware

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

func LanguageMiddleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {

			tr, _ := transport.FromServerContext(ctx)
			lang := tr.RequestHeader().Get("Accept-Language")

			if lang == "" {
				lang = "en"
			}

			ctx = context.WithValue(ctx, "lang", lang)

			return handler(ctx, req)
		}
	}
}