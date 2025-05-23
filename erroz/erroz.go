package erroz

import (
    "errors"
    `fmt`
    
    `dpcms/api/validate`
    "github.com/gin-gonic/gin"
    `github.com/google/uuid`
    `gorm.io/gorm`
)

type businessError struct {
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
    Apply(ctx *gin.Context)
    Abort(ctx *gin.Context)
    ToError() error
    raw() *businessError
}

var _ BusinessError = (*businessError)(nil)

func (s *businessError) Is(err error) bool {
    var bizErr *businessError
    if errors.As(err, &bizErr) {
        return bizErr.Code == s.Code
    }
    return false
}

func (s *businessError) WithOption(option ...Option) BusinessError {
    for _, fn := range option {
        fn(s)
    }
    return s
}

func (s *businessError) Log() BusinessError {
    return &businessError{
        log:    true,
        status: s.status,
        uuid:   s.uuid,
        Debug:  s.Debug,
        Data:   s.Data,
        Code:   s.Code,
        Msg:    s.Msg,
    }
}

func (s *businessError) Format(args ...any) BusinessError {
    return &businessError{
        log:    s.log,
        status: s.status,
        uuid:   s.uuid,
        Debug:  s.Debug,
        Data:   s.Data,
        Code:   s.Code,
        Msg:    fmt.Sprintf(s.Msg, args...),
    }
}

func (s *businessError) Apply(ctx *gin.Context) {
    if s.log {
        errID, err := uuid.NewV7()
        if err != nil {
            _ = ctx.Error(err)
        } else {
            s.uuid = errID.String()
        }
        _ = ctx.Error(s)
    }
    ctx.JSON(s.status, s)
}

func (s *businessError) Abort(ctx *gin.Context) {
    s.Apply(ctx)
    ctx.Abort()
}

func (s *businessError) Error() string {
    if s.uuid == "" {
        return s.Msg
    }
    return s.Msg + "(" + s.uuid + ")"
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
