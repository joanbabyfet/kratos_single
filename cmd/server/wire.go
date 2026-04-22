//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"kratos_single/internal/biz"
	"kratos_single/internal/conf"
	"kratos_single/internal/data"
	"kratos_single/internal/job"
	"kratos_single/internal/pkg/logger"
	"kratos_single/internal/server"
	"kratos_single/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data) (*kratos.App, func(), error) {
	panic(wire.Build(
		logger.ProviderSet, 
		server.ProviderSet, 
		data.ProviderSet, 
		biz.ProviderSet, 
		service.ProviderSet, 
		job.ProviderSet,
		newApp,
	))
}
