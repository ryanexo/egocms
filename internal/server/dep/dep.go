package dep

import (
    `GoBlog/internal/pkg/cache`
    `GoBlog/internal/server/pkg/token`
)

type Dep struct {
    Cache cache.Cache
    Token token.Token
}

func New(c cache.Cache, t token.Token) *Dep {
    return &Dep{
        Cache: c,
        Token: t,
    }
}
