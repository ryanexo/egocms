package middleware

import (
    `errors`
    
    `cms/internal/httpserver`
    `cms/internal/httpserver/validator`
    
    `github.com/gin-gonic/gin`
    `gorm.io/gorm`
)

type ErrorFallback gin.HandlerFunc

func NewErrorFallback() ErrorFallback {
    return func(context *gin.Context) {
        context.Next()
        
        var (
            isResolved      bool
            returnValue     httpserver.Error
            validationError validator.ValidationErrors
        )
        
        err := context.Errors.Last()
        
        switch {
        case errors.As(err, &returnValue):
            break
        
        case errors.As(err, &validationError):
            returnValue = ValidationFailed.WithOption(erroz.WithData(validationError)).prototype()
        
        case errors.Is(err, gorm.ErrRecordNotFound):
            returnValue = DataNotFound.prototype()
        
        default:
            returnValue = Unknown.Wrap(err).prototype()
            isResolved = true
        }
        
        if gin.Mode() == gin.DebugMode && (isResolved || returnValue.Cause != nil) {
            var unResolvedErr error = returnValue
            for {
                e := errors.Unwrap(unResolvedErr)
                if e == nil {
                    break
                }
                returnValue.Debug = append(returnValue.Debug, e.Error())
                unResolvedErr = e
            }
        }
    }
}
