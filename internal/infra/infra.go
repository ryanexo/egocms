package infra

import (
    `dpcms/internal/infra/db`
    `dpcms/internal/infra/hashids`
    `dpcms/internal/infra/logger`
    `dpcms/internal/infra/persistence`
    `dpcms/internal/infra/persistence/query`
    `dpcms/internal/infra/rbac`
    
    `github.com/casbin/casbin/v2`
    "github.com/google/wire"
    "gorm.io/gorm"
)

var InfraProvider = wire.NewSet(
    wire.Struct(new(Infra), "*"),
    logger.New,
    db.NewDB,
    hashids.New,
    persistence.New,
    rbac.New,
)

type Infra struct {
    DB      *gorm.DB
    Query   *query.Query
    Log     *logger.Logger
    HashIds *hashids.HashIds
    Casbin  *casbin.Enforcer
}
