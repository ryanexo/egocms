package service

import (
	"github.com/google/wire"
)

type Services struct {
	User *User
}

var ServiceProviderSet = wire.NewSet(
	wire.Struct(new(Services), "*"),
	NewUserService,
)
