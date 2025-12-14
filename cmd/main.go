package main

import (
    `dpcms/internal/config`
)

var appConfig = config.NewWithBasicConfig()

func main() {
    err := registerCommand()
    if err != nil {
        panic(err)
    }
    launcher, err := createHttpServer(appConfig)
    if err != nil {
        panic(err)
    }
    err = launcher.Run()
    if err != nil {
        panic(err)
    }
}
