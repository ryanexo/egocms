package token

import `github.com/golang-jwt/jwt/v5`

type UserToken struct {
    jwt.RegisteredClaims
    UserID int64 `json:"userID"`
}
