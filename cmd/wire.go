//go:build wireinject

package main

import (
    `dpcms/api/controller`
    `dpcms/api/middleware`
    `dpcms/api/repository`
    `dpcms/api/service`
    `dpcms/config`
    `dpcms/packages`
    `dpcms/server`
    "github.com/google/wire"
)

func createServerLauncher(cfg *config.Config) (*server.Launcher, error) {
    panic(wire.Build(
        config.ProviderSet,
        packages.InfraProviderSet,
        middleware.MiddlewareProviderSet,
        repository.RepoProviderSet,
        service.ServiceProviderSet,
        controller.ControllerProviderSet,
        server.ServerProviderSet,
    ))
}
