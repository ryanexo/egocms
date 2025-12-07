package types

import `time`

type PaginatedResult[T any] struct {
    Total    int64 `json:"total"`
    PageSize int   `json:"pageSize"`
    PageNo   int   `json:"pageNo"`
    List     []T   `json:"list"`
}

type Meta struct {
    ID        int64     `json:"id" validate:"required"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}
