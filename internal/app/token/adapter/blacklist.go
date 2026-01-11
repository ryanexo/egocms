package adapter

import (
    `context`
    `time`
    
    `dpcms/internal/app/token/service`
    `dpcms/internal/infra/persistence/datatype`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/infra/persistence/query`
    
    `gorm.io/gen`
)

type TokenBlacklistRepo struct {
    query *query.Query
}

func NewTokenBlacklistRepo(persist *query.Query) service.TokenBlacklistRepo {
    return &TokenBlacklistRepo{persist}
}

func (r *TokenBlacklistRepo) CloneWithQuery(q *query.Query) service.TokenBlacklistRepo {
    return NewTokenBlacklistRepo(q)
}

func (r *TokenBlacklistRepo) Add(ctx context.Context, userID datatype.SafeUint64, uuid string, expires time.Time) error {
    return r.query.TokenBlacklist.WithContext(ctx).Create(&model.TokenBlacklist{
        UserID:  userID,
        UUID:    uuid,
        Expires: expires,
    })
}

func (r *TokenBlacklistRepo) Remove(ctx context.Context, id datatype.SafeUint64) (gen.ResultInfo, error) {
    dao := r.query.TokenBlacklist
    return dao.WithContext(ctx).Where(dao.ID.Eq(id.Raw())).Delete()
}

func (r *TokenBlacklistRepo) FindByUUID(ctx context.Context, uuid string) (*model.TokenBlacklist, error) {
    dao := r.query.TokenBlacklist
    return dao.WithContext(ctx).Where(dao.UUID.Eq(uuid)).First()
}

func (r *TokenBlacklistRepo) CleanExpired(ctx context.Context, expiredBefore time.Time) (gen.ResultInfo, error) {
    dao := r.query.TokenBlacklist
    return dao.WithContext(ctx).Where(dao.Expires.Lte(expiredBefore)).Delete()
}
