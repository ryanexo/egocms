package service

import (
    `context`
    `errors`
    `time`
    
    `dpcms/api/errors/auth_error`
    `dpcms/api/infra`
    tokenClaim `dpcms/api/service/internal/token`
    `dpcms/config`
    `dpcms/model`
    `dpcms/model/query`
    `github.com/golang-jwt/jwt/v5`
    uuid2 `github.com/google/uuid`
    `gorm.io/gen/field`
    `gorm.io/gorm`
    `gorm.io/gorm/clause`
)

type TokenService struct {
    infra *infra.Infra
    query *query.Query
}

func NewTokenService(infra *infra.Infra) *TokenService {
    return &TokenService{infra: infra, query: query.Use(infra.DB)}
}

func (srv *TokenService) Create(userId int64) (string, error) {
    uuid, err := uuid2.NewV7()
    if err != nil {
        return "", err
    }
    expires := config.Get().Token.Expires
    tokenKey := []byte(config.Get().GlobalKey)
    return jwt.NewWithClaims(jwt.SigningMethodHS512, tokenClaim.UserToken{
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expires))),
            ID:        uuid.String(),
        },
        UserID: userId,
    }).SignedString(tokenKey)
}

func (srv *TokenService) isRevoked(ctx context.Context, userId int64, uuid string, expires int) (bool, error) {
    isRevoked := true
    err := srv.query.Transaction(func(tx *query.Query) error {
        queryCtx := srv.query.WithContext(ctx)
        dao := srv.query.TokenBlacklist
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

func (srv *TokenService) Parse(ctx context.Context, tokenString string) (*model.User, error) {
    claims := &tokenClaim.UserToken{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return config.Get().GlobalKey, nil
    })
    if err != nil || !token.Valid {
        if errors.Is(err, jwt.ErrTokenExpired) {
            return nil, auth_error.ErrAuthorizationExpired.ToError()
        }
        return nil, auth_error.ErrUnauthorized.ToError()
    }
    isRevoked, err := srv.isRevoked(ctx, claims.UserID, claims.ID, config.Get().Token.Expires)
    if err != nil {
        return nil, err
    }
    if isRevoked {
        return nil, auth_error.ErrUnauthorized.ToError()
    }
    return srv.query.User.WithContext(ctx).Where(srv.query.User.ID.Eq(claims.UserID)).First()
}

func (srv *TokenService) Revoke(ctx context.Context, userId int64, uuid string) error {
    return srv.query.WithContext(ctx).TokenBlacklist.Create(&model.TokenBlacklist{
        UserId: userId,
        UUID:   uuid,
    })
}
