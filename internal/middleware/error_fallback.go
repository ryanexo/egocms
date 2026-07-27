package middleware

import (
    `errors`
    
    `cms/internal/httpx`
    `cms/internal/httpx/validator`
    
    `github.com/gin-gonic/gin`
    `gorm.io/gorm`
)

type ErrorFallback gin.HandlerFunc

func NewErrorFallback() ErrorFallback {
    return func(context *gin.Context) {
        context.Next()
        
        err := context.Errors.Last()
        if err == nil {
            return
        }
        
        var response httpx.Error
        var validationError validator.ValidationErrors
        
        switch {
        case errors.As(err, &response):
            response.Abort(context)
        
        case errors.As(err, &validationError):
            httpx.ValidationFailed.Data(validationError).Abort(context)
        
        case errors.Is(err, gorm.ErrRecordNotFound):
            httpx.DataNotFound.Abort(context)
        
        default:
            httpx.Unknown.Abort(context)
        }
    }
}
