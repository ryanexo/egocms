package dto

import (
    "cms/internal/infra/persist/datatype"
    `cms/internal/util/types`
)

type RoleCreateParams struct {
    Name        string                `validate:"required,max=255" json:"name" label:"名称"`
    Description string                `validate:"required,max=255" json:"description" label:"描述"`
    InheritList []datatype.SafeUint64 `validate:"max=10" json:"inheritList" label:"继承角色" swaggertype:"array,string"`
}

type RoleUpdateParams struct {
    types.ResourceID
    RoleCreateParams
}

type Role struct {
    types.Base
    Name        string
    Description string
}

type RoleListParams struct {
    types.Pagination
    Name        *string `json:"name"`
    Description *string `json:"description"`
}

type RoleGrantParams struct {
    types.ResourceID
    PermID []datatype.SafeUint64 `validate:"required" json:"permId" swaggertype:"array,string"`
}
