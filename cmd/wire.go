//go:build wireinject

package main

import (
    `dpcms/internal/config`
    `dpcms/internal/app/controllers`
    `dpcms/internal/app/middleware`
    `dpcms/internal/app/services`
    `dpcms/internal/httpserver`
    `dpcms/internal/infra`
    `dpcms/internal/packages/validate`
    "github.com/google/wire"
)

func createHttpServer(cfg *config.Config) (*httpserver.Launcher, error) {
    panic(wire.Build(
        config.ProviderSet,
        infra.ProviderSet,
        middleware.ProviderSet,
        services.ProviderSet,
        controllers.ProviderSet,
        validate.ProviderSet,
        httpserver.ProviderSet,
    ))
}
