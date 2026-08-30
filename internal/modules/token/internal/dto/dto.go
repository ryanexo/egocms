package dto

import (
    `cms/internal/public/jsontype`
    
    "github.com/golang-jwt/jwt/v5"
)

type UserToken struct {
    jwt.RegisteredClaims
    UserID jsontype.SafeUint64 `json:"uid" apitype:"string"`
}
