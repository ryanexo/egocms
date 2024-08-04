package packages

import (
    `dpcms/packages/cache`
    `dpcms/packages/data`
    `dpcms/packages/logger`
    "github.com/google/wire"
    "go.uber.org/zap"
    "gorm.io/gorm"
)

var InfraProviderSet = wire.NewSet(
    wire.Struct(new(Infra), "*"),
    logger.New,
    data.NewDB,
    cache.New,
)

type Infra struct {
    Cache  *cache.Cache
    DB     *gorm.DB
    Logger *zap.Logger
}
