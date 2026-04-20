package server

import (
	adminv1 "kratos_single/api/admin/v1"
	v1 "kratos_single/api/client/v1"
	commonv1 "kratos_single/api/common/v1"
	"kratos_single/internal/conf"
	"kratos_single/internal/middleware"
	"kratos_single/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, 
	article *service.ArticleService,
	adminArticle *service.AdminArticleService,
	user *service.UserService,
	upload *service.CommonService,
	logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(), //防 panic
			middleware.LanguageMiddleware(), //多语言中间件
			middleware.AuthMiddleware(), //认证中间件
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	// 注册 HTTP 路由 （注册多个服务）
	v1.RegisterArticleHTTPServer(srv, article)
	adminv1.RegisterAdminArticleHTTPServer(srv, adminArticle)
	commonv1.RegisterUploadHTTPServer(srv, upload)
	v1.RegisterUserHTTPServer(srv, user)


	
	return srv
}
