package service

import (
    `dpcms/internal/app/article/service`
    articlemodel `dpcms/internal/app/articlemodel/service`
    category `dpcms/internal/app/category/service`
    menu `dpcms/internal/app/menu/service`
    permission `dpcms/internal/app/permission/service`
    role `dpcms/internal/app/role/service`
    token `dpcms/internal/app/token/service`
    user `dpcms/internal/app/user/service`
    
    "github.com/google/wire"
)

var ServiceProvider = wire.NewSet(
    user.NewUserService,
    category.NewCategoryService,
    token.NewTokenService,
    permission.NewService,
    role.NewRoleService,
    menu.NewMenuService,
    service.NewArticleService,
    articlemodel.NewArticleModelService,
)
