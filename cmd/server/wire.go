//go:build wireinject

package main

import (
    `cms/internal/bootstrap`
    `cms/internal/config`
    
    "github.com/google/wire"
)

func initApp(cfg *config.Config) (bootstrap.Bootstrap, error) {
    panic(wire.Build(
        bootstrap.BootstrapProvider,
    ))
}
