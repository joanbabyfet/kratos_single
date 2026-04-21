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
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, 
	article *service.ArticleService,
	adminArticle *service.AdminArticleService,
	adminAd *service.AdminAdService,
	user *service.UserService,
	upload *service.CommonService,
	ad *service.AdService,
	logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(), //防 panic
			middleware.LanguageMiddleware(), //多语言中间件
			
			//Auth 先执行，写入 role 到 ctx, AdminAuth 再读取 role
			//认证中间件
			selector.Server(middleware.AuthMiddleware()).
			Match(middleware.WhiteListMatcher()). //白名单
			Build(),

			//所有 admin 接口校验 admin 权限
			selector.Server(middleware.AdminAuthMiddleware()).
			Match(middleware.AdminMatcher()).
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
	commonv1.RegisterUploadHTTPServer(srv, upload)
	v1.RegisterUserHTTPServer(srv, user)
	v1.RegisterAdHTTPServer(srv, ad)
	
	return srv
}
