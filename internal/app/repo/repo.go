package repo

import `github.com/google/wire`

type Repo struct {
    Article *Article
}

var RepoProvider = wire.NewSet(
    wire.Struct(new(Repo), "*"),
    NewArticle,
)
