package srvparams

import (
    `dpcms/internal/app/helper/dbscope`
    `dpcms/internal/database/model`
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

type Role struct {
    *model.Role
    InheritList []Role `json:"inheritList"`
}

type RoleListParams struct {
    dbscope.Pagination
    InheritId   *uint64 `json:"inheritId"`
    Name        *string `json:"name"`
    Description *string `json:"description"`
}
