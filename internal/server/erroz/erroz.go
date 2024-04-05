package erroz

import (
    `errors`
    `net/http`

    `github.com/bytedance/sonic`
    `github.com/gin-gonic/gin`
    `github.com/gin-gonic/gin/render`
    `github.com/gookit/validate`
)

type json struct {
    status int
    err    error
    Code   string `json:"code"`
    Msg    string `json:"msg"`
    Data   any    `json:"data,omitempty"`
}

type Option func(*json)

type Response interface {
    render.Render
    WithOption(option ...Option) Response
    Apply(ctx *gin.Context)
    ToError() error
    prototype() *json
}

var _ Response = (*json)(nil)

func (s *json) Render(writer http.ResponseWriter) error {
    data, err := sonic.Marshal(s)
    if err != nil {
        return err
    }
    _, err = writer.Write(data)
    return err
}

func (s *json) WriteContentType(w http.ResponseWriter) {
    w.Header().Set("Content-Type", "application/json")
}

func (s *json) WithOption(option ...Option) Response {
    for _, fn := range option {
        fn(s)
    }
    return s
}

func (s *json) Apply(ctx *gin.Context) {
    ctx.JSON(s.status, s)
    if s.err != nil {
        _ = ctx.Error(s.err)
    }
}

func (s *json) Error() string {
    return s.Msg
}

func (s *json) ToError() error {
    return s
}

func (s *json) prototype() *json {
    return s
}

func Resolve(ctx *gin.Context, err error) {
    var (
        bizErr        *json
        validationErr validate.Errors
        unresolved    error
    )
    switch {
    case errors.As(err, &bizErr):
        break

    case errors.As(err, &validationErr):
        bizErr = ErrValidation.prototype()
        bizErr.WithOption(WithData(validationErr))

    default:
        bizErr = Failed.prototype()
        unresolved = err
    }
    resp := &json{
        err:  unresolved,
        Msg:  bizErr.Msg,
        Code: bizErr.Code,
        Data: bizErr.Data,
    }
    resp.Apply(ctx)
}

func New(code, msg string) Response {
    return &json{Code: code, Msg: msg}
}
