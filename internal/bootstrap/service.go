package bootstrap

import (
    article `dpcms/internal/app/article/service`
    articlemodel `dpcms/internal/app/articlemodel/service`
    category `dpcms/internal/app/category/service`
    config `dpcms/internal/app/config/service`
    fileService `dpcms/internal/app/file/service`
    menu `dpcms/internal/app/menu/service`
    permission `dpcms/internal/app/permission/service`
    role `dpcms/internal/app/role/service`
    token `dpcms/internal/app/token/service`
    user `dpcms/internal/app/user/service`
    
    "github.com/google/wire"
)

var ServiceProvider = wire.NewSet(
    article.NewArticleService,
    articlemodel.NewArticleModelService,
    category.NewCategoryService,
    menu.NewMenuService,
    user.NewUserService,
    role.NewRoleService,
    token.NewTokenService,
    config.NewConfigService,
    fileService.NewFileService,
    wire.Struct(new(permission.PermissionService), "*"),
)
