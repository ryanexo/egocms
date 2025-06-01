package services

import (
    "github.com/google/wire"
)

type Services struct {
    Category CategoryService
    User     UserService
    Token    TokenService
    RBAC     RBACService
    Role     RoleService
    Menu     MenuService
}

var ProviderSet = wire.NewSet(
    wire.Struct(new(Services), "*"),
    NewUserService,
    NewCategoryCategory,
    NewTokenService,
    NewRBACService,
    NewRoleService,
    NewMenuService,
)
