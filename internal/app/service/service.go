package service

import (
    "github.com/google/wire"
)

type Services struct {
    Category     *Category
    User         *User
    Token        *Token
    RBAC         *RBAC
    Role         *Role
    Menu         *Menu
    Article      *Article
    ContentModel *ArticleModel
}

var ServiceProvider = wire.NewSet(
    wire.Struct(new(Services), "*"),
    NewUserService,
    NewCategoryCategory,
    NewTokenService,
    NewRBACService,
    NewRoleService,
    NewMenuService,
    NewArticle,
    NewContentModel,
)
