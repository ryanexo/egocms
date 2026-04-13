package provider

import (
    `cms/internal/httpserver`
    `cms/internal/infra/cache`
    `cms/internal/infra/casbin`
    `cms/internal/infra/db`
    `cms/internal/pkg/hashid`
    `cms/internal/pkg/logger`
    
    `github.com/google/wire`
    
    `cms/internal/infra/persist/contract`
    `cms/internal/infra/persist/query`
    
    `gorm.io/gorm`
)

var InfraProvider = wire.NewSet(
    logger.New,
    db.NewDB,
    casbin.NewRoleCasbin,
    casbin.NewMenuCasbin,
    cache.NewConfigCache,
    hashid.New,
    httpserver.New,
    NewQuery,
    NewTransactor,
)

func NewQuery(db *gorm.DB) *query.Query {
    return query.Use(db)
}

func NewTransactor(q *query.Query) contract.Transactor {
    return q
}
