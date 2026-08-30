package middleware

import (
    `errors`
    `net/http`
    
    `cms/internal/httpx`
    `cms/internal/httpx/validator`
    `cms/internal/public/erroz`
    
    `github.com/gin-gonic/gin`
    `gorm.io/gorm`
)

type ErrorFallback gin.HandlerFunc

func NewErrorFallback() ErrorFallback {
    return func(ctx *gin.Context) {
        ctx.Next()
        
        err := ctx.Errors.Last()
        if err == nil {
            return
        }
        
        var (
            ezErr           *erroz.Error
            validationError validator.ValidationErrors
        )
        
        switch {
        case errors.As(err, &ezErr):
            ctx.AbortWithStatusJSON(http.StatusOK, ezErr.Serialize())
        
        case errors.As(err, &validationError):
            httpx.JSON(ctx, httpx.ValidationFailed.WithData(validationError))
        
        case errors.Is(err, gorm.ErrRecordNotFound):
            httpx.JSON(ctx, httpx.DataNotFound)
        
        default:
            httpx.JSON(ctx, httpx.Unknown)
        }
    }
}
