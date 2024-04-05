//go:build wireinject

package main

import (
    `GoBlog/internal/server`
    `GoBlog/internal/server/biz`
    `GoBlog/internal/server/config`
    `github.com/google/wire`
)

func initServer(c *config.Config) (server.Launcher, error) {
    panic(wire.Build(
        server.ProviderSet,
        biz.ControllerSet,
        biz.ServiceWireSet,
        biz.RepoWireSet,
    ))
}
