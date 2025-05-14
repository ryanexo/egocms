package erroz

import (
    "errors"
    
    `dpcms/api/validate`
    "github.com/gin-gonic/gin"
    `gorm.io/gorm`
)

type businessError struct {
    status int
    Debug  []any  `json:"debug,omitempty"`
    Data   any    `json:"data,omitempty"`
    Code   string `json:"code"`
    Msg    string `json:"msg"`
}

type Option func(*businessError)

type BusinessError interface {
    WithOption(option ...Option) BusinessError
    Apply(ctx *gin.Context)
    Abort(ctx *gin.Context)
    ToError() error
    raw() *businessError
}

var _ BusinessError = (*businessError)(nil)

func (s *businessError) WithOption(option ...Option) BusinessError {
    for _, fn := range option {
        fn(s)
    }
    return s
}

func (s *businessError) Apply(ctx *gin.Context) {
    ctx.JSON(s.status, s)
}

func (s *businessError) Abort(ctx *gin.Context) {
    ctx.AbortWithStatusJSON(s.status, s)
}

func (s *businessError) Error() string {
    return s.Msg
}

func (s *businessError) ToError() error {
    return s
}

func (s *businessError) raw() *businessError {
    return s
}

func Resolve(ctx *gin.Context, err error) {
    var (
        notResolved     bool
        returnValue     *businessError
        validationError validate.ValidationErrors
    )
    
    switch {
    case errors.As(err, &returnValue):
        break
    
    case errors.As(err, &validationError):
        returnValue = ErrValidation.WithOption(WithData(validationError)).raw()
    
    case errors.Is(err, gorm.ErrRecordNotFound):
        returnValue = ErrDataNotFound.raw()
    
    default:
        returnValue = ErrUnknown.raw()
        notResolved = true
        _ = ctx.Error(err)
    }
    
    resp := &businessError{
        Msg:  returnValue.Msg,
        Code: returnValue.Code,
        Data: returnValue.Data,
    }
    if gin.Mode() == gin.DebugMode && notResolved {
        resp.Debug = append(resp.Debug, err.Error())
    }
    resp.Apply(ctx)
}

func ResolveWithAbort(ctx *gin.Context, err error) {
    Resolve(ctx, err)
    ctx.Abort()
}

func New(code, msg string) BusinessError {
    return &businessError{Code: code, Msg: msg}
}
