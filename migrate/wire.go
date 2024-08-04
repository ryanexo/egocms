//go:build wireinject

package main

import (
    `dpcms/packages/data`
    "github.com/google/wire"
    "gorm.io/gorm"
)

func initDB(c *data.DBConfig) (*gorm.DB, error) {
    panic(wire.Build(
        data.NewDB,
    ))
}
