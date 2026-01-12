package erroz

import (
    "errors"
    "fmt"
)

type businessError struct {
    cause  error
    status int
    Result
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

func (s businessError) Write(r Responsible) {
    r.JSON(s.status, s)
}

func (s businessError) WriteWithAbort(r Responsible) {
    s.Write(r)
    r.Abort()
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

func New(code, msg string) BusinessError {
    return &businessError{Result: Result{Code: code, Msg: msg}}
}
