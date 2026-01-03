package erroz

import (
	"errors"
	"fmt"

	"dpcms/internal/httpserver/validator"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type businessError struct {
	cause  error
	status int
	Result
}

type Option func(*businessError)

type BusinessError interface {
	WithOption(option ...Option) BusinessError
	Format(...any) BusinessError
	Wrap(err error) BusinessError
	Write(ctx *gin.Context)
	WriteWithAbort(ctx *gin.Context)
	ToError() error
	prototype() businessError
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
	s.cause = err
	return s
}

func (s businessError) Unwrap() error {
	return s.cause
}

func (s businessError) WithOption(option ...Option) BusinessError {
	for _, fn := range option {
		fn(&s)
	}
	return s
}

func (s businessError) Format(args ...any) BusinessError {
	s.Msg = fmt.Sprintf(s.Msg, args...)
	return s
}

func (s businessError) Write(ctx *gin.Context) {
	ctx.JSON(s.status, s)
}

func (s businessError) WriteWithAbort(ctx *gin.Context) {
	s.Write(ctx)
	ctx.Abort()
}

func (s businessError) Error() string {
	var errMsg string
	if s.cause != nil {
		errMsg = s.cause.Error()
	} else {
		errMsg = s.Msg
	}
	return errMsg
}

func (s businessError) ToError() error {
	return s
}

func (s businessError) prototype() businessError {
	return s
}

func ResolveWithWrite(ctx *gin.Context, err error) {
	var (
		notResolved     bool
		returnValue     businessError
		validationError validator.ValidationErrors
	)

	switch {
	case errors.As(err, &returnValue):
		break

	case errors.As(err, &validationError):
		returnValue = ValidationFailed.WithOption(WithData(validationError)).prototype()

	case errors.Is(err, gorm.ErrRecordNotFound):
		returnValue = DataNotFound.prototype()

	default:
		returnValue = Unknown.Wrap(err).prototype()
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
	return &businessError{Result: Result{Code: code, Msg: msg}}
}
