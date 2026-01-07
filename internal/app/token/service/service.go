package service

import (
    `context`
    `errors`
    `time`
    
    `dpcms/internal/app/token/internal/dto`
    tokenBlacklistRepo `dpcms/internal/app/token/repo`
    userRepo `dpcms/internal/app/user/repo`
    `dpcms/internal/config`
    `dpcms/internal/erroz`
    `dpcms/internal/infra`
    `dpcms/internal/infra/persistence/datatype`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/infra/persistence/query`
    
    `github.com/golang-jwt/jwt/v5`
    `github.com/google/uuid`
    `gorm.io/gorm`
)

type TokenService struct {
    config  *config.Config
    persist *query.Query
}

func NewTokenService(config *config.Config, i *infra.Infra) *TokenService {
    return &TokenService{persist: i.Query, config: config}
}

func (s TokenService) Create(userId datatype.SafeUint64) (string, error) {
    uuidValue, err := uuid.NewV7()
    if err != nil {
        return "", err
    }
    expires := s.config.Token.Expires
    tokenKey := []byte(s.config.GlobalKey)
    return jwt.NewWithClaims(
        jwt.SigningMethodHS512,
        dto.UserToken{
            RegisteredClaims: jwt.RegisteredClaims{
                ExpiresAt: jwt.NewNumericDate(
                    time.Now().Add(time.Duration(expires)),
                ),
                ID: uuidValue.String(),
            },
            UserID: userId,
        }).SignedString(tokenKey)
}

func (s TokenService) isRevoked(ctx context.Context, uuid string, expires int) (bool, error) {
    tbRepo := tokenBlacklistRepo.NewTokenBlacklistRepository(s.persist)
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

func (s TokenService) Parse(ctx context.Context, tokenString string) (*dto.UserToken, error) {
    claims := &dto.UserToken{}
    jwtData, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return s.config.GlobalKey, nil
    })
    if err != nil || !jwtData.Valid {
        if errors.Is(err, jwt.ErrTokenExpired) {
            return nil, erroz.AuthorizationExpired.ToError()
        }
        return nil, erroz.Unauthorized.Wrap(err).ToError()
    }
    isRevoked, err := s.isRevoked(ctx, claims.ID, s.config.Token.Expires)
    if err != nil {
        return nil, err
    }
    if isRevoked {
        return nil, erroz.Unauthorized.ToError()
    }
    return claims, nil
}

func (s TokenService) GetUserFromToken(ctx context.Context, tokenString string) (*model.User, error) {
    claims, err := s.Parse(ctx, tokenString)
    if err != nil {
        return nil, err
    }
    return userRepo.NewUserRepo(s.persist).FindByID(ctx, claims.UserID)
}

func (s TokenService) Revoke(ctx context.Context, tokenString string) error {
    claims, err := s.Parse(ctx, tokenString)
    if err != nil {
        return err
    }
    return s.persist.WithContext(ctx).TokenBlacklist.Create(&model.TokenBlacklist{
        UserId: claims.UserID,
        UUID:   claims.ID,
    })
}
