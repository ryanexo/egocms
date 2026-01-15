package controller

import (
    `reflect`
    
    article `cms/internal/app/article/controller`
    articleModel `cms/internal/app/articlemodel/controller`
    category `cms/internal/app/category/controller`
    file `cms/internal/app/file/controller`
    menu `cms/internal/app/menu/controller`
    permission `cms/internal/app/permission/controller`
    role `cms/internal/app/role/controller`
    user `cms/internal/app/user/controller`
    `cms/internal/httpserver`
    `cms/internal/util/reflectutil`
    
    `github.com/google/wire`
)

var ControllerProvider = wire.NewSet(
    wire.Struct(new(ControllerSet), "*"),
    article.NewArticleController,
    articleModel.NewArticleModelController,
    category.NewCategoryController,
    menu.NewMenuController,
    role.NewRoleController,
    user.NewUserController,
    file.NewFileController,
    permission.NewPermissionController,
    NewRouteRegistrar,
)

type ControllerSet struct {
    User         *user.UserController
    Menu         *menu.MenuController
    Category     *category.CategoryController
    Role         *role.RoleController
    Article      *article.ArticleController
    ArticleModel *articleModel.ArticleModelController
    File         *file.FileController
    Perm         *permission.PermissionController
}

var _ httpserver.Route = (*ControllerSet)(nil)

func (c ControllerSet) Setup(router httpserver.Router) {
    _ = reflectutil.InvokeImplementedStruct[httpserver.Route](c, func(_ reflect.Value, i httpserver.Route) error {
        i.Setup(router)
        return nil
    })
}

func NewRouteRegistrar(ctlSet *ControllerSet) httpserver.Route {
    return ctlSet
}
