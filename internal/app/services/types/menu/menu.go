package menu

import (
    `dpcms/internal/app/services/types`
    `dpcms/internal/utils/dbscopes`
)

type RetrieveListParams struct {
    dbscopes.Pagination
    Name      *string `json:"name"`
    Ancestor  *int64  `json:"ancestor"`
    Recursive bool    `json:"recursive"`
}

type DetailParams struct {
    ID int64 `validate:"required" json:"id"`
}

type CreateParams struct {
    ParentID int64  `validate:"required" json:"parentId"`
    Name     string `validate:"required" json:"name"`
    Sequence int64  `json:"sequence"`
    URI      string `validate:"required,alphanum" json:"uri"`
    Template string `json:"template"`
    Remark   string `json:"remark"`
}

type UpdateParams struct {
    CreateParams
    ID int64 `validate:"required" json:"id"`
}

type DeleteParams struct {
    ID int64 `validate:"required" json:"id"`
}

type Detail struct {
    types.Meta
    ParentID int64  `json:"parentID"`
    Name     string `json:"name"`
    Sequence int64  `json:"sequence"`
    Visible  bool   `json:"visible"`
    URI      string `json:"URI"`
    Template string `json:"template"`
    Remark   string `json:"remark"`
}
