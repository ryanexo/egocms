package srvparams

import (
    `dpcms/internal/app/helper/dbscope`
    `dpcms/internal/packages/database`
)

type RoleCreateParams struct {
    Name        string   `validate:"required;max=255" json:"name" label:"名称"`
    Description string   `validate:"required;max=255" json:"description" label:"描述"`
    InheritList []uint64 `validate:"max=10" json:"inheritList" label:"继承角色"`
}

type RoleUpdateParams struct {
    RoleCreateParams
    ID uint64 `validate:"required" json:"id"`
}

type DeleteParams struct {
    ID uint64 `validate:"required" json:"id"`
}

type Detail struct {
    database.Model
    Name        string   `json:"name"`
    Description string   `json:"description"`
    InheritList []Detail `json:"inheritList"`
}

type ListDetail struct {
    ID          uint64 `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
}

type ListRetrieveParams struct {
    dbscope.Pagination
    Name        *string `json:"name"`
    Description *string `json:"description"`
}
