package erroz

type Result struct {
    Debug []any  `json:"debug,omitempty" swaggerignore:"true"`
    Data  any    `json:"data,omitempty"`
    Code  string `json:"code"`
    Msg   string `json:"msg"`
}
