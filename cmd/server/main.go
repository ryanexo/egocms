package main

import (
    `fmt`
    `runtime/debug`
    
    _ `dpcms/docs`
    `dpcms/internal/config`
)

// @title           EgoCms Api
// @version         0.1
// @license.name    Apache 2.0
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @BasePath        /
func main() {
    appConfig := &config.Config{}
    if err := registerCommand(appConfig); err != nil {
        panic(fmt.Sprintf("%+v\n%s", err, debug.Stack()))
    }
    err := appConfig.Validate()
    if err != nil {
        panic(fmt.Sprintf("%+v\n%s", err, debug.Stack()))
    }
    launcher, err := createHttpServer(appConfig)
    if err != nil {
        panic(fmt.Sprintf("%+v\n%s", err, debug.Stack()))
    }
    err = launcher.Run(true)
    if err != nil {
        panic(fmt.Sprintf("%+v\n%s", err, debug.Stack()))
    }
}
