package infra

import (
    `dpcms/internal/packages/database`
    `dpcms/internal/packages/logger`
    
    "github.com/google/wire"
    "gorm.io/gorm"
)

var ProviderSet = wire.NewSet(
    wire.Struct(new(Infra), "*"),
    logger.New,
    database.NewDB,
)

type Infra struct {
    DB  *gorm.DB
    Log *logger.Logger
}
