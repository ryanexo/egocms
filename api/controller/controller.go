package controller

import (
    "reflect"
    
    `dpcms/server`
    "github.com/gin-gonic/gin"
    `github.com/google/wire`
)

var ControllerProviderSet = wire.NewSet(
    wire.Struct(new(Controller), "*"),
    NewUserController,
    NewRouteRegistrar,
)

type Controller struct {
    User *UserController
}

type IController interface {
    setup(engine *gin.Engine)
}

var _ server.Routes = (*Controller)(nil)

func (c Controller) SetupRoutes(engine *gin.Engine) {
    ref := reflect.ValueOf(c)
    
    for i := 0; i < ref.NumField(); i++ {
        iterateField := ref.Field(i)
        if iterateField.Kind() == reflect.Ptr {
            iterateField = iterateField.Elem()
        }
        if iterateField.Kind() != reflect.Struct {
            continue
        }
        if !iterateField.CanInterface() {
            continue
        }
        iterateController, ok := iterateField.Interface().(IController)
        if !ok || iterateController == nil {
            continue
        }
        iterateController.setup(engine)
    }
}

func NewRouteRegistrar(c *Controller) server.Routes {
    return c
}
