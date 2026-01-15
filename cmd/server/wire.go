//go:build wireinject

package main

import (
    `cms/internal/bootstrap`
    config2 `cms/internal/bootstrap/config`
    `cms/internal/bootstrap/controller`
    `cms/internal/bootstrap/filedriver`
    `cms/internal/bootstrap/infra`
    `cms/internal/bootstrap/middleware`
    `cms/internal/bootstrap/repo`
    `cms/internal/bootstrap/service`
    `cms/internal/config`
    
    "github.com/google/wire"
)

func initApp(cfg *config.Config) (bootstrap.Bootstrap, error) {
    panic(wire.Build(
        infra.InfraProvider,
        middleware.MiddlewareProvider,
        repo.RepoProvider,
        service.ServiceProvider,
        controller.ControllerProvider,
        filedriver.FileDriverProvider,
        config2.ConfigProvider,
        bootstrap.BootstrapProvider,
    ))
}
