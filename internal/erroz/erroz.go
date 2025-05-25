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
    Write(ctx *gin.Context)
    WriteWithAbort(ctx *gin.Context)
    ToError() error
    raw() *businessError
}

var _ BusinessError = (*businessError)(nil)

func (s *businessError) clone() *businessError {
    return &businessError{
        log:    s.log,
        status: s.status,
        uuid:   s.uuid,
        Debug:  s.Debug,
        Data:   s.Data,
        Code:   s.Code,
        Msg:    s.Msg,
    }
}

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
    err := s.clone()
    err.log = true
    return err
}

func (s *businessError) Format(args ...any) BusinessError {
    err := s.clone()
    err.Msg = fmt.Sprintf(s.Msg, args...)
    return err
}

func (s *businessError) Write(ctx *gin.Context) {
    if s.log {
        bizErr := s
        errID, err := uuid.NewV7()
        if err != nil {
            _ = ctx.Error(err)
        } else {
            bizErr = s.clone()
            bizErr.uuid = errID.String()
        }
        _ = ctx.Error(bizErr)
    }
    ctx.JSON(s.status, s)
}

func (s *businessError) WriteWithAbort(ctx *gin.Context) {
    s.Write(ctx)
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

func ResolveWithWrite(ctx *gin.Context, err error) {
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
