package httpserver

import (
    "fmt"
    `net/http`
    
    `github.com/gin-gonic/gin`
)

type errorMeta struct {
    Message string `json:"message"`
    Code    int    `json:"code"`
    Debug   any    `json:"debug,omitempty"`
    Data    any    `json:"data"`
}

type Error struct {
    status int
    meta   errorMeta
}

func (e Error) Is(err error) bool {
    t, ok := err.(Error)
    if !ok {
        return false
    }
    
    return err == e || e.meta.Code == t.meta.Code
}

func (e Error) Data(data any) Error {
    e.meta.Data = data
    return e
}

func (e Error) Debug(data any) Error {
    if gin.Mode() == gin.DebugMode {
        e.meta.Debug = data
    }
    return e
}

func (e Error) Format(args ...any) Error {
    e.meta.Message = fmt.Sprintf(e.meta.Message, args...)
    return e
}

func (e Error) Error() string {
    return e.meta.Message
}

func (e Error) JSON(ctx *gin.Context) {
    ctx.JSON(e.status, e)
}

func (e Error) Abort(ctx *gin.Context) {
    ctx.JSON(e.status, e)
    ctx.Abort()
}

func NewError(code int, cause string, status ...int) *Error {
    httpStatus := http.StatusOK
    if len(status) > 0 {
        httpStatus = status[0]
    }
    return &Error{status: httpStatus, meta: errorMeta{Message: cause, Code: code}}
}
