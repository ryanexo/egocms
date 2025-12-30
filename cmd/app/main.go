package main

import (
    `dpcms/cmd/app`
    `dpcms/internal/config`
)

var appConfig = config.NewWithBasicConfig()

func main() {
    err := app.registerCommand()
    if err != nil {
        panic(err)
    }
    launcher, err := app.createHttpServer(appConfig)
    if err != nil {
        panic(err)
    }
    err = launcher.Run()
    if err != nil {
        panic(err)
    }
}
