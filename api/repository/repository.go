package repository

import "github.com/google/wire"

type Repositories struct {
    User     *User
    Category *Category
    Token    *Token
}

var ProviderSet = wire.NewSet(
    wire.Struct(new(Repositories), "*"),
    NewUserRepo,
    NewCategoryRepo,
    NewTokenRepo,
)
