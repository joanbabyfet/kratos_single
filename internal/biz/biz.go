package biz

import "github.com/google/wire"

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewArticleUsecase, 
	NewUserUsecase, 
	NewAdUsecase, 
	NewAdminUsecase, 
	NewMailUsecase, 
	NewRBACUsecase,
	wire.Bind(new(RBAC), new(*RBACUsecase)), //加这个否则报错
)
