package authz

import (
    `dpcms/internal/erroz`
    `dpcms/internal/constant`
    
    `github.com/gin-gonic/gin`
)

func GetCurrentUser[T any](ctx *gin.Context) (T, error) {
    var zero T
    val, ok := ctx.Get(constant.RequestUserKey)
    if !ok {
        return zero, erroz.Unauthorized.ToError()
    }
    user, ok := val.(T)
    if !ok {
        return zero, erroz.Unauthorized.ToError()
    }
    return user, nil
}
