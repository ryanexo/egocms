package dto

import (
    "cms/internal/infra/persistence/datatype"
    `cms/internal/types`
)

type RoleCreateParams struct {
    Name        string                `validate:"required,max=255" json:"name" label:"名称"`
    Description string                `validate:"required,max=255" json:"description" label:"描述"`
    InheritList []datatype.SafeUint64 `validate:"max=10" json:"inheritList" label:"继承角色" swaggertype:"array,string"`
}

type RoleUpdateParams struct {
    RoleCreateParams
    ID datatype.SafeUint64 `validate:"required" json:"id" swaggertype:"string"`
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
