package service

import (
    "github.com/google/wire"
)

type Services struct {
    User  *User
    Token *Token
}

var ProviderSet = wire.NewSet(
    wire.Struct(new(Services), "*"),
    NewUserService,
    NewTokenService,
)
