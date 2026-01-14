package bootstrap

import (
    `cms/internal/httpserver`
    `cms/internal/infra/cache`
    `cms/internal/infra/casbin`
    `cms/internal/infra/db`
    `cms/internal/infra/logger`
    `cms/internal/infra/persistence`
    `cms/internal/infra/xhashids`
    
    `github.com/google/wire`
)

var BootstrapProvider = wire.NewSet(
    logger.New,
    db.NewDB,
    persistence.NewQuery,
    persistence.NewTxManager,
    casbin.NewRoleCasbin,
    casbin.NewMenuCasbin,
    cache.NewConfigCache,
    xhashids.New,
    httpserver.New,
)
