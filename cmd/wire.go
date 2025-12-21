//go:build wireinject

package main

import (
    `dpcms/internal/app/controller`
    `dpcms/internal/app/middleware`
    `dpcms/internal/app/repo`
    `dpcms/internal/app/service`
    `dpcms/internal/config`
    `dpcms/internal/httpserver`
    `dpcms/internal/infra`
    
    "github.com/google/wire"
)

func createHttpServer(cfg *config.Config) (*httpserver.Launcher, error) {
    panic(wire.Build(
        config.ProviderSet,
        infra.ProviderSet,
        middleware.ProviderSet,
        service.ProviderSet,
        controller.ProviderSet,
        httpserver.ProviderSet,
        repo.ProviderSet,
    ))
}
