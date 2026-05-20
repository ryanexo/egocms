package erroz

import (
    `errors`
    
    `cms/internal/httpserver`
    `cms/internal/httpserver/erroz`
    
    `github.com/gin-gonic/gin`
    `github.com/go-playground/validator/v10`
    `gorm.io/gorm`
)

func ResolveWithWrite(ctx erroz.Responsible, err error) {
    var (
        notResolved     bool
        returnValue     httpserver.Error
        validationError validator.ValidationErrors
    )
    
    switch {
    case errors.As(err, &returnValue):
        break
    
    case errors.As(err, &validationError):
        returnValue = ValidationFailed.WithOption(erroz.WithData(validationError)).prototype()
    
    case errors.Is(err, gorm.ErrRecordNotFound):
        returnValue = DataNotFound.prototype()
    
    default:
        returnValue = Unknown.Wrap(err).prototype()
        notResolved = true
    }
    
    if gin.Mode() == gin.DebugMode && (notResolved || returnValue.Cause != nil) {
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
    returnValue.Write(ctx)
}

func ResolveWithAbort(ctx erroz.Responsible, err error) {
    ResolveWithWrite(ctx, err)
    ctx.Abort()
}
