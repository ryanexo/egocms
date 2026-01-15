package infra

import (
    `cms/internal/httpserver`
    `cms/internal/infra/cache`
    `cms/internal/infra/casbin`
    `cms/internal/infra/db`
    `cms/internal/infra/logger`
    `cms/internal/infra/persist`
    `cms/internal/infra/xhashids`
    
    `github.com/google/wire`
)

var InfraProvider = wire.NewSet(
    logger.New,
    db.NewDB,
    persist.NewQuery,
    persist.NewTxManager,
    casbin.NewRoleCasbin,
    casbin.NewMenuCasbin,
    cache.NewConfigCache,
    xhashids.New,
    httpserver.New,
)
