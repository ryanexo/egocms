package swaggertype

// PaginatedResult[T]
type PaginatedResult[T any] = Result[T]

// Result[T]
type Result[T any] struct {
    Data T      `json:"data,omitempty"`
    Code string `json:"code"`
    Msg  string `json:"msg"`
}
