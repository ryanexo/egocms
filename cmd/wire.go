//go:build wireinject

package main

import (
    `dpcms/internal/config`
    `dpcms/internal/app/controller`
    `dpcms/internal/app/middleware`
    `dpcms/internal/app/service`
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
        service.ProviderSet,
        controller.ProviderSet,
        validate.ProviderSet,
        httpserver.ProviderSet,
    ))
}
