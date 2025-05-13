//go:build wireinject

package main

import (
    `dpcms/api/controller`
    `dpcms/api/infra`
    `dpcms/api/middleware`
    `dpcms/api/repository`
    `dpcms/api/service`
    `dpcms/api/validate`
    `dpcms/config`
    `dpcms/server`
    "github.com/google/wire"
)

func createServerLauncher(cfg *config.Config) (*server.Launcher, error) {
    panic(wire.Build(
        config.ProviderSet,
        infra.ProviderSet,
        middleware.ProviderSet,
        repository.ProviderSet,
        service.ProviderSet,
        controller.ProviderSet,
        validate.ProviderSet,
        server.ProviderSet,
    ))
}
