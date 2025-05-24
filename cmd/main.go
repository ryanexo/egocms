package main

import (
    `dpcms/config`
)

func main() {
    err := registerCommand()
    if err != nil {
        panic(err)
    }
    cfg := config.NewWithDefaultConfig()
    launcher, err := createHttpServer(cfg)
    if err != nil {
        panic(err)
    }
    err = launcher.Run()
    if err != nil {
        panic(err)
    }
}
