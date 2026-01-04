package dtotype

// PaginatedResult[T]
type PaginatedResult[T any] struct {
    Total    int64 `json:"total"`
    PageSize int   `json:"pageSize"`
    PageNo   int   `json:"pageNo"`
    List     []T   `json:"list" swaggerignore:"true"`
}
