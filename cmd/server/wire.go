//go:build wireinject

package main

import (
    `dpcms/internal/bootstrap`
    `dpcms/internal/config`
    `dpcms/internal/httpserver`
    
    "github.com/google/wire"
)

func createHttpServer(cfg *config.Config) (*httpserver.Launcher, error) {
    panic(wire.Build(
        config.ConfigProvider,
        bootstrap.BootstrapProvider,
        bootstrap.MiddlewareProvider,
        bootstrap.RepoProvider,
        bootstrap.ServiceProvider,
        bootstrap.ControllerProvider,
        bootstrap.FileDriverProvider,
    ))
}
