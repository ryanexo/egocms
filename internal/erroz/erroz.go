package erroz

import (
    "errors"
    `fmt`
    
    `dpcms/internal/packages/validate`
    "github.com/gin-gonic/gin"
    `github.com/google/uuid`
    `gorm.io/gorm`
)

type businessError struct {
    err    error
    log    bool
    status int
    uuid   string
    Debug  []any  `json:"debug,omitempty"`
    Data   any    `json:"data,omitempty"`
    Code   string `json:"code"`
    Msg    string `json:"msg"`
}

type Option func(*businessError)

type BusinessError interface {
    WithOption(option ...Option) BusinessError
    Log() BusinessError
    Format(...any) BusinessError
    Wrap(err error) BusinessError
    Write(ctx *gin.Context)
    WriteWithAbort(ctx *gin.Context)
    ToError() error
    raw() businessError
}

var _ BusinessError = (*businessError)(nil)

func (s businessError) Is(err error) bool {
    var bizErr businessError
    if errors.As(err, &bizErr) {
        return bizErr.Code == s.Code
    }
    return false
}

func (s businessError) Wrap(err error) BusinessError {
    s.err = err
    return s
}

func (s businessError) Unwrap() error {
    return s.err
}

func (s businessError) WithOption(option ...Option) BusinessError {
    for _, fn := range option {
        fn(&s)
    }
    return s
}

func (s businessError) Log() BusinessError {
    s.log = true
    return s
}

func (s businessError) Format(args ...any) BusinessError {
    s.Msg = fmt.Sprintf(s.Msg, args...)
    return s
}

func (s businessError) Write(ctx *gin.Context) {
    if s.log {
        id, err := uuid.NewV7()
        if err != nil {
            _ = ctx.Error(err)
        } else {
            idStr := id.String()
            s.uuid = idStr
            s.Msg = s.Msg + "[" + idStr + "]"
        }
        _ = ctx.Error(s)
    }
    ctx.JSON(s.status, s)
}

func (s businessError) WriteWithAbort(ctx *gin.Context) {
    s.Write(ctx)
    ctx.Abort()
}

func (s businessError) Error() string {
    var errMsg string
    if s.err != nil {
        errMsg = s.err.Error()
    } else {
        errMsg = s.Msg
    }
    if s.uuid == "" {
        return errMsg
    }
    return "[" + s.uuid + "]" + errMsg
}

func (s businessError) ToError() error {
    return s
}

func (s businessError) raw() businessError {
    return s
}

func ResolveWithWrite(ctx *gin.Context, err error) {
    var (
        notResolved     bool
        returnValue     businessError
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
        returnValue = ErrUnknown.Log().Wrap(err).raw()
        notResolved = true
    }
    
    if gin.Mode() == gin.DebugMode && notResolved {
        returnValue.Debug = append(returnValue.Debug, err.Error())
    }
    returnValue.Write(ctx)
}

func ResolveWithAbort(ctx *gin.Context, err error) {
    ResolveWithWrite(ctx, err)
    ctx.Abort()
}

func New(code, msg string) BusinessError {
    return &businessError{Code: code, Msg: msg}
}
