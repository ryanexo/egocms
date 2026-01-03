package token

import (
	"dpcms/internal/infra/persistence/datatype"

	"github.com/golang-jwt/jwt/v5"
)

type UserToken struct {
	jwt.RegisteredClaims
	UserID datatype.SafeUint64 `json:"uid" swaggertype:"string"`
}
