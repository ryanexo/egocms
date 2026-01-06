package service

import (
    `context`
    `errors`
    `time`
    
    `dpcms/internal/app/erroz`
    `dpcms/internal/app/service/internal/token`
    `dpcms/internal/config`
    `dpcms/internal/infra`
    `dpcms/internal/infra/persistence/datatype`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/infra/persistence/query`
    `dpcms/internal/infra/persistence/repo`
    
    `github.com/golang-jwt/jwt/v5`
    `github.com/google/uuid`
    `gorm.io/gorm`
)

type Token struct {
    config  *config.Config
    persist *query.Query
}

func NewTokenService(config *config.Config, i *infra.Infra) *Token {
    return &Token{persist: i.Query, config: config}
}

func (srv Token) Create(userId datatype.SafeUint64) (string, error) {
    uuidValue, err := uuid.NewV7()
    if err != nil {
        return "", err
    }
    expires := srv.config.Token.Expires
    tokenKey := []byte(srv.config.GlobalKey)
    return jwt.NewWithClaims(
        jwt.SigningMethodHS512,
        token.UserToken{
            RegisteredClaims: jwt.RegisteredClaims{
                ExpiresAt: jwt.NewNumericDate(
                    time.Now().Add(time.Duration(expires)),
                ),
                ID: uuidValue.String(),
            },
            UserID: userId,
        }).SignedString(tokenKey)
}

func (srv Token) isRevoked(ctx context.Context, uuid string, expires int) (bool, error) {
    tbRepo := repo.NewTokenBlacklistRepo(srv.persist)
    data, err := tbRepo.FindByUUID(ctx, uuid)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return false, nil
        }
        return false, err
    }
    
    expiresTime := data.CreatedAt.Add(time.Duration(expires) * time.Second)
    if expiresTime.Before(time.Now()) {
        return true, nil
    }
    _, err = tbRepo.Remove(ctx, data.ID)
    return true, err
}

func (srv Token) Parse(ctx context.Context, tokenString string) (*token.UserToken, error) {
    claims := &token.UserToken{}
    jwtData, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return srv.config.GlobalKey, nil
    })
    if err != nil || !jwtData.Valid {
        if errors.Is(err, jwt.ErrTokenExpired) {
            return nil, erroz.AuthorizationExpired.ToError()
        }
        return nil, erroz.Unauthorized.Wrap(err).ToError()
    }
    isRevoked, err := srv.isRevoked(ctx, claims.ID, srv.config.Token.Expires)
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
    return repo.NewUserRepo(srv.persist).FindByID(ctx, claims.UserID)
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
