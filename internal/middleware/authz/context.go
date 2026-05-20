package authz

import (
    `context`
    
    `github.com/gin-gonic/gin`
)

type contextKeyType struct{}

var contextKey = contextKeyType{}

func setCurrentUser(ctx *gin.Context, user any) {
    ctx.Set(contextKey, user)
}

func GetCurrentUser[T any](ctx context.Context) (T, error) {
    var zero T
    val := ctx.Value(contextKey)
    if val == nil {
        return zero, ErrAuthorized
    }
    user, ok := val.(T)
    if !ok {
        return zero, ErrAuthorized
    }
    return user, nil
}
