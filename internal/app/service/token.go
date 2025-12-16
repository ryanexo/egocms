package service

import (
    `context`
    `errors`
    `time`
    
    `dpcms/internal/app/erroz`
    tokenClaim `dpcms/internal/app/service/internal/token`
    `dpcms/internal/config`
    `dpcms/internal/database/model`
    `dpcms/internal/database/query`
    `dpcms/internal/infra`
    
    `github.com/golang-jwt/jwt/v5`
    uuid2 `github.com/google/uuid`
    `gorm.io/gen/field`
    `gorm.io/gorm`
    `gorm.io/gorm/clause`
)

type TokenService struct {
    config *config.Config
    infra  *infra.Infra
    query  *query.Query
}

func NewTokenService(config *config.Config, infra *infra.Infra) *TokenService {
    return &TokenService{infra: infra, query: query.Use(infra.DB), config: config}
}

func (srv TokenService) Create(userId uint64) (string, error) {
    uuid, err := uuid2.NewV7()
    if err != nil {
        return "", err
    }
    expires := srv.config.Token.Expires
    tokenKey := []byte(srv.config.GlobalKey)
    return jwt.NewWithClaims(jwt.SigningMethodHS512, tokenClaim.UserToken{
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expires))),
            ID:        uuid.String(),
        },
        UserID: userId,
    }).SignedString(tokenKey)
}

func (srv TokenService) isRevoked(ctx context.Context, userId uint64, uuid string, expires int) (bool, error) {
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

func (srv TokenService) Parse(ctx context.Context, tokenString string) (*tokenClaim.UserToken, error) {
    claims := &tokenClaim.UserToken{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return srv.config.GlobalKey, nil
    })
    if err != nil || !token.Valid {
        if errors.Is(err, jwt.ErrTokenExpired) {
            return nil, erroz.ErrAuthorizationExpired.ToError()
        }
        return nil, erroz.ErrUnauthorized.ToError()
    }
    isRevoked, err := srv.isRevoked(ctx, claims.UserID, claims.ID, srv.config.Token.Expires)
    if err != nil {
        return nil, err
    }
    if isRevoked {
        return nil, erroz.ErrUnauthorized.ToError()
    }
    return claims, nil
}

func (srv TokenService) GetUserFromToken(ctx context.Context, tokenString string) (*model.User, error) {
    claims, err := srv.Parse(ctx, tokenString)
    if err != nil {
        return nil, err
    }
    uo := srv.query.User
    return uo.WithContext(ctx).Preload(uo.Profile).Where(uo.ID.Eq(claims.UserID)).First()
}

func (srv TokenService) Revoke(ctx context.Context, tokenString string) error {
    claims, err := srv.Parse(ctx, tokenString)
    if err != nil {
        return err
    }
    return srv.query.WithContext(ctx).TokenBlacklist.Create(&model.TokenBlacklist{
        UserId: claims.UserID,
        UUID:   claims.ID,
    })
}
