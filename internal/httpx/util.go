package httpx

import (
    `net/http`
    
    `cms/internal/public/erroz`
    
    `github.com/gin-gonic/gin`
)

func HandlerFunc(fn func(ctx *gin.Context) error) gin.HandlerFunc {
    return func(ctx *gin.Context) {
        if err := fn(ctx); err != nil {
            _ = ctx.Error(err)
        }
    }
}

func JSON(ctx *gin.Context, res *erroz.Error) {
    ctx.AbortWithStatusJSON(http.StatusOK, res.Serialize())
}
