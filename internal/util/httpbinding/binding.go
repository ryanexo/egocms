package httpbinding

import (
    `cms/internal/erroz`
    erroz2 `cms/internal/httpserver/erroz`
    
    `github.com/gin-gonic/gin`
)

func BindJSON[T any](ctx *gin.Context, fn func(params T) (any, error)) {
    var params T
    if err := ctx.ShouldBindJSON(&params); err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    
    result, err := fn(params)
    if err != nil {
        erroz.ResolveWithWrite(ctx, err)
    } else if result == nil {
        erroz.OK.Write(ctx)
    } else {
        erroz.OK.WithOption(erroz2.WithData(result)).Write(ctx)
    }
}
