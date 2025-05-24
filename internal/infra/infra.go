package infra

import (
    `dpcms/internal/packages/cache`
    `dpcms/internal/packages/database`
    `dpcms/internal/packages/logger`
    "github.com/google/wire"
    "go.uber.org/zap"
    "gorm.io/gorm"
)

var ProviderSet = wire.NewSet(
    wire.Struct(new(Infra), "*"),
    logger.New,
    database.NewDB,
    cache.New,
)

type Infra struct {
    Cache  *cache.Cache
    DB     *gorm.DB
    Logger *zap.Logger
}
