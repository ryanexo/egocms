package dto

import (
    "cms/internal/infra/persistence/datatype"
    `cms/internal/types`
)

type Menu struct {
    types.Base
    ParentID datatype.SafeUint64 `json:"parentId"`
    Type     int8                `json:"type"`
    Name     string              `json:"name"`
    Sequence datatype.SafeInt64  `json:"sequence"`
    Visible  *datatype.BoolInt8  `json:"visible,omitempty" swaggertype:"boolean"`
    URI      string              `json:"uri"`
    Template string              `json:"template"`
    Remark   string              `json:"remark"`
}

type MenuListQueryParams struct {
    types.Pagination
    Name     *string              `json:"name"`
    ParentID *datatype.SafeUint64 `json:"parentId" swaggertype:"string"`
}

type MenuCreateParams struct {
    ParentID datatype.SafeUint64 `json:"parentId" swaggertype:"string"`
    Type     int8                `json:"type"`
    Name     string              `validate:"required" json:"name" label:"名称"`
    Sequence datatype.SafeInt64  `json:"sequence" swaggertype:"string"`
    URI      string              `validate:"required,alphanum" json:"uri"`
    Template string              `json:"template"`
    Remark   string              `json:"remark"`
}

type MenuUpdateParams struct {
    ID       datatype.SafeUint64 `validate:"required" json:"id" swaggertype:"string"`
    Type     int8                `json:"type"`
    Name     string              `json:"name"`
    Sequence datatype.SafeInt64  `json:"sequence" swaggertype:"string"`
    URI      string              `json:"uri"`
    Template string              `json:"template"`
    Remark   string              `json:"remark"`
}

type MenuMoveParams struct {
    ID       datatype.SafeUint64 `validate:"required" json:"id" swaggertype:"string"`
    TargetID datatype.SafeUint64 `validate:"required" json:"targetId" swaggertype:"string"`
}
