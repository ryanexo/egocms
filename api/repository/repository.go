package repository

import "github.com/google/wire"

type Repositories struct {
    User     *User
    Category *Category
}

var RepoProviderSet = wire.NewSet(
    wire.Struct(new(Repositories), "*"),
    NewUserRepo,
    NewCategory,
)
