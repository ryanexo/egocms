package main

import (
    _ `dpcms/docs`
    `dpcms/internal/config`
)

// @title           EgoCms Api
// @version         1.0
// @license.name    Apache 2.0

// @BasePath        /
func main() {
    appConfig := &config.Config{}
    if err := registerCommand(appConfig); err != nil {
        panic(err)
    }
    launcher, err := createHttpServer(appConfig)
    if err != nil {
        panic(err)
    }
    err = launcher.Run(true)
    if err != nil {
        panic(err)
    }
}
