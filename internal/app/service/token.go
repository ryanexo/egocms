package service

import (
    `context`
    `errors`
    `time`
    
    `dpcms/internal/app/erroz`
    `dpcms/internal/app/service/internal/token`
    `dpcms/internal/config`
    `dpcms/internal/infra`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/infra/persistence/query`
    
    `github.com/golang-jwt/jwt/v5`
    `github.com/google/uuid`
    `gorm.io/gen/field`
    `gorm.io/gorm`
    `gorm.io/gorm/clause`
)

type Token struct {
    config  *config.Config
    persist *query.Query
}

func NewTokenService(config *config.Config, i *infra.Infra) *Token {
    return &Token{persist: i.Query, config: config}
}

func (srv Token) Create(userId uint64) (string, error) {
    uuid, err := uuid.NewV7()
    if err != nil {
        return "", err
    }
    expires := srv.config.Token.Expires
    tokenKey := []byte(srv.config.GlobalKey)
    return jwt.NewWithClaims(jwt.SigningMethodHS512, token.UserToken{
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expires))),
            ID:        uuid.String(),
        },
        UserID: userId,
    }).SignedString(tokenKey)
}

func (srv Token) isRevoked(ctx context.Context, userId uint64, uuid string, expires int) (bool, error) {
    isRevoked := true
    err := srv.persist.Transaction(func(tx *query.Query) error {
        queryCtx := srv.persist.WithContext(ctx)
        dao := srv.persist.TokenBlacklist
        _, err := queryCtx.TokenBlacklist.Clauses(clause.Locking{Strength: "UPDATE"}).Where(dao.UserId.Eq(userId)).Select(field.NewUnsafeFieldRaw("1")).Find()
        if err != nil {
            return err
        }
        token, err := queryCtx.TokenBlacklist.Where(dao.UserId.Eq(userId), dao.UUID.Eq(uuid)).First()
        if err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                isRevoked = false
                return nil
            }
            return err
        }
        exp := token.CreatedAt.Add(time.Duration(expires) * time.Second)
        if exp.Before(time.Now()) {
            return nil
        }
        _, err = queryCtx.TokenBlacklist.Delete(token)
        if err == nil {
            isRevoked = false
        }
        return err
    })
    if err != nil {
        return isRevoked, err
    }
    return isRevoked, nil
}

func (srv Token) Parse(ctx context.Context, tokenString string) (*token.UserToken, error) {
    claims := &token.UserToken{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return srv.config.GlobalKey, nil
    })
    if err != nil || !token.Valid {
        if errors.Is(err, jwt.ErrTokenExpired) {
            return nil, erroz.AuthorizationExpired.ToError()
        }
        return nil, erroz.Unauthorized.ToError()
    }
    isRevoked, err := srv.isRevoked(ctx, claims.UserID, claims.ID, srv.config.Token.Expires)
    if err != nil {
        return nil, err
    }
    if isRevoked {
        return nil, erroz.Unauthorized.ToError()
    }
    return claims, nil
}

func (srv Token) GetUserFromToken(ctx context.Context, tokenString string) (*model.User, error) {
    claims, err := srv.Parse(ctx, tokenString)
    if err != nil {
        return nil, err
    }
    uo := srv.persist.User
    return uo.WithContext(ctx).Preload(uo.Profile).Where(uo.ID.Eq(claims.UserID)).First()
}

func (srv Token) Revoke(ctx context.Context, tokenString string) error {
    claims, err := srv.Parse(ctx, tokenString)
    if err != nil {
        return err
    }
    return srv.persist.WithContext(ctx).TokenBlacklist.Create(&model.TokenBlacklist{
        UserId: claims.UserID,
        UUID:   claims.ID,
    })
}
