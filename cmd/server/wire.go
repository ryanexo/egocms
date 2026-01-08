//go:build wireinject

package main

import (
    `dpcms/internal/bootstrap`
    `dpcms/internal/bootstrap/controller`
    `dpcms/internal/bootstrap/middleware`
    `dpcms/internal/bootstrap/service`
    `dpcms/internal/config`
    `dpcms/internal/httpserver`
    
    "github.com/google/wire"
)

func createHttpServer(cfg *config.Config) (*httpserver.Launcher, error) {
    panic(wire.Build(
        config.ConfigProvider,
        bootstrap.BootstrapProvider,
        middleware.MiddlewareProvider,
        controller.ControllerProvider,
        service.ServiceProvider,
    ))
}
