package contract

import (
    "context"
    "time"
    
    `cms/internal/public/jsontype`
    
    "cms/internal/infra/store/model"
    
    "gorm.io/gen"
)

type TokenBlacklistRepo interface {
    persistence.Repository[TokenBlacklistRepo]
    Add(ctx context.Context, userID jsontype.SafeUint64, uuid string, expires time.Time) error
    Remove(ctx context.Context, id jsontype.SafeUint64) (gen.ResultInfo, error)
    FindByUUID(ctx context.Context, uuid string) (*model.TokenBlacklist, error)
    CleanExpired(ctx context.Context, expiredBefore time.Time) (gen.ResultInfo, error)
}
