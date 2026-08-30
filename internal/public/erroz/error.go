package erroz

import `fmt`

type Error struct {
    Message string
    Code    string
    Data    any
    Cause   error
}

func (e *Error) clone(patch func(*Error)) *Error {
    clone := *e
    patch(&clone)
    return &clone
}

func (e *Error) Serialize() map[string]any {
    result := map[string]any{
        "message": e.Message,
        "code":    e.Code,
    }
    if e.Data != nil {
        result["data"] = e.Data
    }
    return result
}

func (e *Error) Is(err error) bool {
    t, ok := err.(*Error)
    if !ok {
        return false
    }
    
    return e.Code == t.Code
}

func (e *Error) Unwrap() error {
    return e.Cause
}

func (e *Error) WithData(data any) *Error {
    return e.clone(func(e *Error) {
        e.Data = data
    })
}

func (e *Error) Format(args ...any) *Error {
    return e.clone(func(e *Error) {
        e.Message = fmt.Sprintf(e.Message, args...)
    })
}

func (e *Error) Wrap(err error) *Error {
    e.Cause = err
    return e
}

func (e *Error) Error() string {
    return e.Message
}

func NewError(code string, message string) *Error {
    return &Error{Code: code, Message: message}
}
