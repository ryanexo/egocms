package dto

import (
    `dpcms/internal/infra/persistence/datatype`
    `dpcms/internal/types`
)

type Permission struct {
    types.Base
    MenuID      *datatype.SafeUint64 `json:"menuId" swaggertype:"string"`
    Name        string               `validate:"required" json:"name" label:"权限名称"`
    Description string               `json:"description"`
    Resource    string               `validate:"required" json:"resource" label:"资源标识"`
    Action      string               `validate:"required" json:"action" label:"操作标识"`
}

type PermissionCreateParams struct {
    MenuID      *datatype.SafeUint64 `json:"menuId" swaggertype:"string"`
    Name        string               `json:"name"`
    Description string               `json:"description"`
    Resource    string               `json:"resource"`
    Action      string               `json:"action"`
}

type PermissionUpdateParams struct {
    ID datatype.SafeUint64 `json:"id" swaggertype:"string"`
    PermissionCreateParams
}
