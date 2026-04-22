package server

import (
	adminv1 "kratos_single/api/admin/v1"
	v1 "kratos_single/api/client/v1"
	"kratos_single/internal/conf"
	"kratos_single/internal/middleware"
	"kratos_single/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, 
	article *service.ArticleService,
	adminArticle *service.AdminArticleService,
	adminAd *service.AdminAdService,
	user *service.UserService,
	ad *service.AdService,
	admin *service.AdminService,
	logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(), //防 panic
			middleware.LanguageMiddleware(), //多语言中间件
			middleware.RateLimit(), //限流
			
			// ===============================
			// 1. 非白名单接口先做登录认证
			// 白名单接口自动跳过
			// ===============================
			selector.Server(middleware.AuthMiddleware()).
				Match(middleware.WhiteListMatcher()).
				Build(),

			// ===============================
			// 2. 后台接口：admin 权限校验
			// /v1/admin/*
			// ===============================
			selector.Server(middleware.AdminAuthMiddleware()).
				Match(middleware.AdminMatcher()).
				Build(),

			// ===============================
			// 3. 前台接口：普通用户登录校验（可选）
			// /v1/api/*
			// ===============================
			selector.Server(middleware.ClientAuthMiddleware()).
				Match(middleware.ClientMatcher()).
				Build(),
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
	adminv1.RegisterAdminAdHTTPServer(srv, adminAd)
	adminv1.RegisterAdminHTTPServer(srv, admin)
	v1.RegisterUserHTTPServer(srv, user)
	v1.RegisterAdHTTPServer(srv, ad)
	
	// upload 路由（走原生）
	srv.HandleFunc("/v1/upload", UploadHandler)

	return srv
}
