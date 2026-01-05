package dto

import (
    `dpcms/internal/infra/persistence/datatype`
    `dpcms/internal/infra/persistence/dbscope`
)

type Category struct {
    Base
    ParentID datatype.SafeUint64 `json:"parentId"`
    Sequence datatype.SafeInt64  `json:"sequence"`
    Name     string              `json:"name"`
    Path     string              `json:"path"`
    Type     int8                `json:"type"`
    Visible  *datatype.BoolInt8  `json:"visible,omitempty" swaggertype:"bool"`
    SEO      *CategorySEO        `json:"seo"`
}

type CategorySEO struct {
    Title       string `json:"title"`
    Keywords    string `json:"keywords"`
    Description string `json:"description"`
}

type CategoryCreateParams struct {
    ParentID datatype.SafeUint64 `json:"parentID" swaggertype:"string"`
    Sequence datatype.SafeInt64  `json:"sequence"`
    Type     int8                `json:"type"`
    Name     string              `validate:"required" json:"name" label:"名称"`
    Path     string              `validate:"required" json:"path" label:"路径"`
    Visible  *datatype.BoolInt8  `json:"visible" swaggertype:"bool"`
    SEO      CategorySEO         `json:"seo"`
}

type CategoryUpdateParams struct {
    ID       datatype.SafeUint64 `validate:"required" json:"id" swaggertype:"string"`
    Sequence datatype.SafeInt64  `json:"sequence"`
    Type     int8                `json:"type"`
    Name     string              `json:"name"`
    Path     string              `json:"path"`
    Visible  *datatype.BoolInt8  `json:"visible" swaggertype:"bool"`
    SEO      CategorySEO         `json:"seo"`
}

type CategoryMoveParams struct {
    ID       datatype.SafeUint64 `json:"id" swaggertype:"string"`
    TargetID datatype.SafeUint64 `json:"targetId" swaggertype:"string"`
}

type CategoryListParams struct {
    dbscope.Pagination
    ID       *datatype.SafeUint64 `json:"id" swaggertype:"string"`
    ParentID *datatype.SafeUint64 `json:"parentId" swaggertype:"string"`
    Type     *int8                `json:"type"`
    Name     *string              `json:"name"`
    Path     *string              `json:"path"`
    Visible  *datatype.BoolInt8   `json:"visible" swaggertype:"bool"`
}
