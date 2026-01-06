package swaggertype

type Result[T any] struct {
    Data T      `json:"data,omitempty"`
    Code string `json:"code"`
    Msg  string `json:"msg"`
}

type PaginatedResult[T any] = Result[[]T]

type CreateResult struct {
    Data string `json:"data" example:"123456"` // 数据 ID
    Code string `json:"code"`
    Msg  string `json:"msg"`
}

type EmptyResult struct {
    Code string `json:"code"`
    Msg  string `json:"msg"`
}
