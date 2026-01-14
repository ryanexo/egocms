//go:build wireinject

package main

import (
    `cms/internal/bootstrap`
    `cms/internal/config`
    `cms/internal/httpserver`
    
    "github.com/google/wire"
)

func createHttpServer(cfg *config.Config) (*httpserver.Launcher, error) {
    panic(wire.Build(
        bootstrap.BootstrapProvider,
        bootstrap.MiddlewareProvider,
        bootstrap.RepoProvider,
        bootstrap.ServiceProvider,
        bootstrap.ControllerProvider,
        bootstrap.FileDriverProvider,
        bootstrap.ConfigProvider,
    ))
}
