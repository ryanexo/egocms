package dto

import (
	"cms/internal/pkg/datatype"

	"github.com/golang-jwt/jwt/v5"
)

type UserToken struct {
	jwt.RegisteredClaims
	UserID datatype.SafeUint64 `json:"uid" swaggertype:"string"`
}
