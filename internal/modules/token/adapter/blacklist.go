package adapter

import (
    "context"
    "time"
    
    `cms/internal/infra/store/gorm/gquery`
    `cms/internal/public/jsontype`
    
    "cms/internal/infra/store/model"
    "cms/internal/modules/token/contract"
    
    "gorm.io/gen"
)

type TokenBlacklistRepo struct {
    query *gquery.Query
}

func NewTokenBlacklistRepo(persist *gquery.Query) contract.TokenBlacklistRepo {
    return &TokenBlacklistRepo{persist}
}

func (r *TokenBlacklistRepo) CloneWithQuery(q *gquery.Query) contract.TokenBlacklistRepo {
    return NewTokenBlacklistRepo(q)
}

func (r *TokenBlacklistRepo) Add(ctx context.Context, userID jsontype.SafeUint64, uuid string, expires time.Time) error {
    return r.query.TokenBlacklist.WithContext(ctx).Create(&model.TokenBlacklist{
        UserID:  userID,
        UUID:    uuid,
        Expires: expires,
    })
}

func (r *TokenBlacklistRepo) Remove(ctx context.Context, id jsontype.SafeUint64) (gen.ResultInfo, error) {
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
