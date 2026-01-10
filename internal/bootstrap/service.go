package bootstrap

import (
    article `dpcms/internal/app/article/service`
    articlemodel `dpcms/internal/app/articlemodel/service`
    category `dpcms/internal/app/category/service`
    config `dpcms/internal/app/config/service`
    menu `dpcms/internal/app/menu/service`
    permission `dpcms/internal/app/permission/service`
    role `dpcms/internal/app/role/service`
    token `dpcms/internal/app/token/service`
    user `dpcms/internal/app/user/service`
    
    "github.com/google/wire"
)

var ServiceProvider = wire.NewSet(
    wire.Struct(new(user.UserService), "*"),
    wire.Struct(new(category.CategoryService), "*"),
    wire.Struct(new(token.TokenService), "*"),
    wire.Struct(new(permission.PermissionService), "*"),
    wire.Struct(new(role.RoleService), "*"),
    wire.Struct(new(menu.MenuService), "*"),
    wire.Struct(new(article.ArticleService), "*"),
    wire.Struct(new(articlemodel.ArticleModelService), "*"),
    wire.Struct(new(config.ConfigService), "*"),
)
