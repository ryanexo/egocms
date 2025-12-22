package dto

import (
    `dpcms/internal/infra/persistence/dbscope`
)

type MenuListQueryParams struct {
    dbscope.Pagination
    Name     *string `json:"name"`
    Ancestor *uint64 `json:"ancestor"`
}

type MenuCreateParams struct {
    ParentID uint64 `json:"parentId"`
    Name     string `validate:"required" json:"name" label:"名称"`
    Sequence int64  `json:"sequence"`
    URI      string `validate:"required,alphanum" json:"uri"`
    Template string `json:"template"`
    Remark   string `json:"remark"`
}

type MenuUpdateParams struct {
    ID       uint64 `validate:"required" json:"id"`
    Name     string `json:"name"`
    Sequence int64  `json:"sequence"`
    URI      string `json:"uri"`
    Template string `json:"template"`
    Remark   string `json:"remark"`
}

type MenuMoveParams struct {
    ID       uint64 `validate:"required" json:"id"`
    TargetID uint64 `validate:"required" json:"targetId"`
}
