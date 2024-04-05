package token

import (
    `time`

    `GoBlog/internal/server/config`
    `github.com/golang-jwt/jwt/v5`
)

type Token struct {
    key []byte
}

func (t Token) Key() []byte {
    return t.key
}

func New(cfg *config.Config) Token {
    return Token{key: []byte(cfg.Server.SecureKey)}
}

func (t Token) CreateUserToken(claims UserTokenClaims) (string, error) {
    if claims.ExpiresAt.IsZero() {
        claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 3))
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(t.key)
}
