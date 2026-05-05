package service

import (
    `context`
    `errors`
    `time`
    
    `cms/internal/app/token/contract`
    `cms/internal/app/token/internal/dto`
    `cms/internal/app/token/internal/errno`
    user `cms/internal/app/user/contract`
    `cms/internal/config`
    `cms/internal/infra/persistence/datatype`
    `cms/internal/infra/persistence/model`
    
    `github.com/golang-jwt/jwt/v5`
    `github.com/google/uuid`
    `gorm.io/gorm`
)

type TokenService struct {
    config        *config.Config
    blacklistRepo contract.TokenBlacklistRepo
    userRepo      user.UserRepo
}

func NewTokenService(config *config.Config, blacklistRepo contract.TokenBlacklistRepo, userRepo user.UserRepo) *TokenService {
    return &TokenService{
        config:        config,
        blacklistRepo: blacklistRepo,
        userRepo:      userRepo,
    }
}

func (s TokenService) Create(userID datatype.SafeUint64) (string, error) {
    uuidValue, err := uuid.NewV7()
    if err != nil {
        return "", err
    }
    tokenKey := []byte(s.config.AppKey)
    return jwt.NewWithClaims(
        jwt.SigningMethodHS512,
        dto.UserToken{
            RegisteredClaims: jwt.RegisteredClaims{
                ExpiresAt: jwt.NewNumericDate(
                    time.Now().Add(time.Hour * 24 * 7),
                ),
                ID: uuidValue.String(),
            },
            UserID: userID,
        }).SignedString(tokenKey)
}

func (s TokenService) isRevoked(ctx context.Context, uuid string) (bool, error) {
    data, err := s.blacklistRepo.FindByUUID(ctx, uuid)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return false, nil
        }
        return false, err
    }
    
    if data.Expires.Before(time.Now()) {
        return true, nil
    }
    _, err = s.blacklistRepo.Remove(ctx, data.ID)
    return true, err
}

func (s TokenService) Parse(ctx context.Context, tokenString string) (*dto.UserToken, error) {
    claims := &dto.UserToken{}
    jwtData, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return s.config.AppKey, nil
    })
    if err != nil || !jwtData.Valid {
        if errors.Is(err, jwt.ErrTokenExpired) {
            return nil, errno.AuthorizationExpired.ToError()
        }
        return nil, errno.Unauthorized.Wrap(err).ToError()
    }
    isRevoked, err := s.isRevoked(ctx, claims.ID)
    if err != nil {
        return nil, err
    }
    if isRevoked {
        return nil, errno.Unauthorized.ToError()
    }
    return claims, nil
}

func (s TokenService) GetUserFromToken(ctx context.Context, tokenString string) (*model.User, error) {
    claims, err := s.Parse(ctx, tokenString)
    if err != nil {
        return nil, err
    }
    return s.userRepo.FindByID(ctx, claims.UserID)
}

func (s TokenService) Revoke(ctx context.Context, tokenString string) error {
    claims, err := s.Parse(ctx, tokenString)
    if err != nil {
        return err
    }
    return s.blacklistRepo.Add(ctx, claims.UserID, claims.ID, claims.ExpiresAt.Time)
}
