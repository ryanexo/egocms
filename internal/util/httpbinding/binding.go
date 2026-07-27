package httpbinding

import (
    `cms/internal/httpx`
    
    `github.com/gin-gonic/gin`
)

func BindJSON[T any](ctx *gin.Context, fn func(params T) (any, error)) error {
    var params T
    if err := ctx.ShouldBindJSON(&params); err != nil {
        return err
    }
    
    result, err := fn(params)
    if err != nil {
        return err
    } else if result == nil {
        httpx.OK.JSON(ctx)
    } else {
        httpx.OK.Data(result).JSON(ctx)
    }
    return nil
}
