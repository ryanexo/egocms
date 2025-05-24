//go:build wireinject

package main

import (
    `dpcms/config`
    `dpcms/internal/http/controller`
    `dpcms/internal/http/middleware`
    `dpcms/internal/http/service`
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
