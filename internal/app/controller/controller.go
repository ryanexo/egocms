package controller

import (
    "reflect"
    
    `dpcms/internal/httpserver`
    
    "github.com/gin-gonic/gin"
    `github.com/google/wire`
)

var ProviderSet = wire.NewSet(
    wire.Struct(new(Controllers), "*"),
    NewRouteRegistrar,
    NewUserController,
    NewMenuController,
)

type Controllers struct {
    User UserController
    Menu MenuController
}

type IController interface {
    setup(engine *gin.Engine)
}

var _ httpserver.Routes = (*Controllers)(nil)

func (c Controllers) SetupRoutes(engine *gin.Engine) {
    ref := reflect.ValueOf(c)
    
    for i := 0; i < ref.NumField(); i++ {
        iterateField := ref.Field(i)
        if !iterateField.CanInterface() {
            continue
        }
    SETUP:
        controller, ok := iterateField.Interface().(IController)
        if ok {
            controller.setup(engine)
        } else if iterateField.Kind() == reflect.Ptr {
            iterateField = iterateField.Elem()
            goto SETUP
        }
    }
}

func NewRouteRegistrar(c *Controllers) httpserver.Routes {
    return c
}
