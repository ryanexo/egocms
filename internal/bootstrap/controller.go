package bootstrap

import (
    "reflect"
    
    article `dpcms/internal/app/article/controller`
    articleModel `dpcms/internal/app/articlemodel/controller`
    category `dpcms/internal/app/category/controller`
    menu `dpcms/internal/app/menu/controller`
    role `dpcms/internal/app/role/controller`
    user `dpcms/internal/app/user/controller`
    `dpcms/internal/httpserver`
    
    "github.com/gin-gonic/gin"
    `github.com/google/wire`
)

var ControllerProvider = wire.NewSet(
    wire.Struct(new(Controllers), "*"),
    wire.Struct(new(article.ArticleController), "*"),
    wire.Struct(new(articleModel.ArticleModelController), "*"),
    wire.Struct(new(category.CategoryController), "*"),
    wire.Struct(new(menu.MenuController), "*"),
    wire.Struct(new(role.RoleController), "*"),
    wire.Struct(new(user.UserController), "*"),
    NewRouteRegistrar,
)

type Controllers struct {
    User         user.UserController
    Menu         menu.MenuController
    Category     category.CategoryController
    Role         role.RoleController
    Article      article.ArticleController
    ArticleModel articleModel.ArticleModelController
}

type IController interface {
    Setup(engine *gin.Engine)
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
            controller.Setup(engine)
        } else if iterateField.Kind() == reflect.Ptr {
            iterateField = iterateField.Elem()
            goto SETUP
        }
    }
}

func NewRouteRegistrar(c *Controllers) httpserver.Routes {
    return c
}
