package service

import (
    "github.com/google/wire"
)

type Services struct {
    User  *UserService
    Token *TokenService
}

var ProviderSet = wire.NewSet(
    wire.Struct(new(Services), "*"),
    NewUserService,
    NewCategory,
    NewTokenService,
)
