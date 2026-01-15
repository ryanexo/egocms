package types

import (
    `time`
    
    `cms/internal/infra/persist/datatype`
)

type ResourceID struct {
    ID datatype.SafeUint64 `validate:"required" json:"id" swaggertype:"string" example:"123456"`
}

type Base struct {
    ID        datatype.SafeUint64 `json:"id" swaggertype:"string"`
    CreatedAt time.Time           `json:"createdAt"`
    UpdatedAt time.Time           `json:"updatedAt"`
}

type PaginatedResult[T any] struct {
    Pagination
    Total    int64 `json:"total"`
    PageSize int   `json:"pageSize"`
    PageNo   int   `json:"pageNo"`
    List     []T   `json:"list" swaggerignore:"true"`
}

type Pagination struct {
    PageNo   int `json:"pageNo"`
    PageSize int `json:"pageSize"`
}
