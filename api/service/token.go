package service

import (
    `context`
    `errors`
    `time`
    
    `dpcms/api/infra`
    `dpcms/api/repository`
    tokenClaim `dpcms/api/service/internal/token`
    `dpcms/config`
    `dpcms/erroz`
    `dpcms/model`
    `github.com/golang-jwt/jwt/v5`
    uuid2 `github.com/google/uuid`
)

type Token struct {
    infra infra.Infra
    repo  repository.Repositories
}

func NewTokenService(infra infra.Infra, repo repository.Repositories) *Token {
    return &Token{infra, repo}
}

func (t *Token) Create(userId uint) (string, error) {
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

func (t *Token) Parse(tokenString string) (*model.User, error) {
    claims := &tokenClaim.UserToken{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return config.Get().GlobalKey, nil
    })
    if err != nil || !token.Valid {
        if errors.Is(err, jwt.ErrTokenExpired) {
            return nil, erroz.ErrAuthorizationExpired.ToError()
        }
        return nil, erroz.ErrUnauthorized.ToError()
    }
    isRevoked, err := t.repo.Token.IsRevoked(context.Background(), claims.UserID, claims.ID, config.Get().Token.Expires)
    if err != nil {
        return nil, err
    }
    if isRevoked {
        return nil, erroz.ErrUnauthorized.ToError()
    }
    return t.repo.User.FindByID(context.Background(), claims.UserID)
}
