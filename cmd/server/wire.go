//go:build wireinject

package main

import (
    `dpcms/internal/config`
    `dpcms/internal/httpserver`
    `dpcms/internal/infra`
    `dpcms/internal/provider/controller`
    `dpcms/internal/provider/middleware`
    `dpcms/internal/provider/service`
    
    "github.com/google/wire"
)

func createHttpServer(cfg *config.Config) (*httpserver.Launcher, error) {
    panic(wire.Build(
        config.ConfigProvider,
        infra.InfraProvider,
        controller.ControllerProvider,
        service.ServiceProvider,
        middleware.MiddlewareProvider,
        httpserver.ProviderSet,
    ))
}
