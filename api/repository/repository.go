package repository

import "github.com/google/wire"

type Repositories struct {
    Category *Category
    Token    *Token
}

var ProviderSet = wire.NewSet(
    wire.Struct(new(Repositories), "*"),
    NewCategoryRepo,
    NewTokenRepo,
)
