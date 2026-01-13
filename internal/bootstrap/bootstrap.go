package bootstrap

import (
    `dpcms/internal/httpserver`
    `dpcms/internal/infra/cache`
    `dpcms/internal/infra/db`
    `dpcms/internal/infra/logger`
    `dpcms/internal/infra/persistence`
    `dpcms/internal/infra/rbac`
    
    `github.com/google/wire`
)

var BootstrapProvider = wire.NewSet(
    httpserver.New,
    rbac.NewRoleCasbin,
    rbac.NewMenuCasbin,
    logger.New,
    db.NewDB,
    persistence.NewQuery,
    persistence.NewTxManager,
    cache.NewConfigCache,
    NewXHashIds,
)
