package main

import (
    _ "embed"
    
    `dpcms/config`
)

func main() {
    err := registerCommand()
    if err != nil {
        panic(err)
    }
    cfg := config.NewWithDefaultConfig()
    launcher, err := createServerLauncher(cfg)
    if err != nil {
        panic(err)
    }
    err = launcher.Run()
    if err != nil {
        panic(err)
    }
}
