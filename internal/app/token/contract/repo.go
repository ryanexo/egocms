package contract

import (
	"cms/internal/infra/persistence"
	"cms/internal/pkg/datatype"
	"context"
	"time"

	"cms/internal/infra/persistence/model"

	"gorm.io/gen"
)

type TokenBlacklistRepo interface {
	persistence.Repository[TokenBlacklistRepo]
	Add(ctx context.Context, userID datatype.SafeUint64, uuid string, expires time.Time) error
	Remove(ctx context.Context, id datatype.SafeUint64) (gen.ResultInfo, error)
	FindByUUID(ctx context.Context, uuid string) (*model.TokenBlacklist, error)
	CleanExpired(ctx context.Context, expiredBefore time.Time) (gen.ResultInfo, error)
}
