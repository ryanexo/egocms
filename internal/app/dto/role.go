package dto

import (
    `dpcms/internal/persistence/dbscope`
    `dpcms/internal/persistence/model`
)

type RoleCreateParams struct {
    Name        string   `validator:"required;max=255" json:"name" label:"名称"`
    Description string   `validator:"required;max=255" json:"description" label:"描述"`
    InheritList []uint64 `validator:"max=10" json:"inheritList" label:"继承角色"`
}

type RoleUpdateParams struct {
    RoleCreateParams
    ID uint64 `validator:"required" json:"id"`
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
