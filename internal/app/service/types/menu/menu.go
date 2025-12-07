package menu

import (
    `dpcms/internal/app/helper/dbscope`
    `dpcms/internal/app/service/types`
)

type ListParams struct {
    dbscope.Pagination
    Name      *string `json:"name"`
    Ancestor  *int64  `json:"ancestor"`
    Recursive bool    `json:"recursive"`
}

type CreateParams struct {
    ParentID int64  `json:"parentId"`
    Name     string `validate:"required" json:"name" label:"名称"`
    Sequence int64  `json:"sequence"`
    URI      string `validate:"required,alphanum" json:"uri"`
    Template string `json:"template"`
    Remark   string `json:"remark"`
}

type UpdateParams struct {
    ID       int64  `validate:"required" json:"id"`
    ParentID int64  `json:"parentId"`
    Name     string `json:"name"`
    Sequence int64  `json:"sequence"`
    URI      string `json:"uri"`
    Template string `json:"template"`
    Remark   string `json:"remark"`
}

type DetailResult struct {
    types.Meta
    CreateParams
    Visible bool `json:"visible"`
}
