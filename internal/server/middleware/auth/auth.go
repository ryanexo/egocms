package auth

import (
    `GoBlog/internal/server/config`
    `GoBlog/internal/server/enum`
    `GoBlog/internal/server/erroz`
    `github.com/gin-gonic/gin`
    `github.com/golang-jwt/jwt/v5`
)

func New(cfg *config.Config) gin.HandlerFunc {
    jwtKey := []byte(cfg.Server.SecureKey)

    return func(ctx *gin.Context) {
        credential := ctx.GetHeader("Authorization")

        token, parseErr := jwt.Parse(
            credential,
            func(token *jwt.Token) (interface{}, error) {
                return jwtKey, nil
            },
        )

        if parseErr != nil || token.Valid {
            erroz.ErrUnauthorized.Apply(ctx)
            return
        }

        ctx.Set(enum.MiddlewareAuth, token.Claims)
        ctx.Next()
    }
}
