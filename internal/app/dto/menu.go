package dto

import (
    `dpcms/internal/infra/persistence/dbscope`
    `dpcms/internal/infra/persistence/datatype`
)

type MenuListQueryParams struct {
    dbscope.Pagination
    Name     *string              `json:"name"`
    Ancestor *datatype.SafeUint64 `json:"ancestor" swaggertype:"string"`
}

type MenuCreateParams struct {
    ParentID datatype.SafeUint64 `json:"parentId" swaggertype:"string"`
    Name     string              `validate:"required" json:"name" label:"名称"`
    Sequence datatype.SafeInt64  `json:"sequence" swaggertype:"string"`
    URI      string              `validate:"required,alphanum" json:"uri"`
    Template string              `json:"template"`
    Remark   string              `json:"remark"`
}

type MenuUpdateParams struct {
    ID       datatype.SafeUint64 `validate:"required" json:"id" swaggertype:"string"`
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
