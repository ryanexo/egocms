package service

import (
    `context`
    `time`
    
    `cms/internal/infra/persistence/contract`
    `cms/internal/infra/persistence/datatype`
    `cms/internal/infra/persistence/model`
    
    `gorm.io/gen`
)

type TokenBlacklistRepo interface {
    contract.Repository[TokenBlacklistRepo]
    Add(ctx context.Context, userID datatype.SafeUint64, uuid string, expires time.Time) error
    Remove(ctx context.Context, id datatype.SafeUint64) (gen.ResultInfo, error)
    FindByUUID(ctx context.Context, uuid string) (*model.TokenBlacklist, error)
    CleanExpired(ctx context.Context, expiredBefore time.Time) (gen.ResultInfo, error)
}
