package middleware

import (
    `errors`
    
    `cms/internal/erroz`
    `cms/internal/httpserver`
    `cms/internal/httpserver/validator`
    
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
        
        var response httpserver.Error
        var validationError validator.ValidationErrors
        
        switch {
        case errors.As(err, &response):
            response.Abort(context)
        
        case errors.As(err, &validationError):
            erroz.ValidationFailed.Data(validationError).Abort(context)
        
        case errors.Is(err, gorm.ErrRecordNotFound):
            erroz.DataNotFound.Abort(context)
        
        default:
            erroz.Unknown.Abort(context)
        }
    }
}
