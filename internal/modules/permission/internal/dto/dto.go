package dto

import (
    `cms/internal/public/apitype`
    `cms/internal/public/jsontype`
)

type Permission struct {
    apitype.Base
    MenuName    string               `json:"menuName"`
    MenuID      *jsontype.SafeUint64 `json:"menuId" apitype:"string"`
    Name        string               `validate:"required" json:"name" label:"权限名称"`
    Description string               `json:"description"`
    Resource    string               `validate:"required" json:"resource" label:"资源标识"`
    Action      string               `validate:"required" json:"action" label:"操作标识"`
}

type PermissionCreateParams struct {
    MenuID      *jsontype.SafeUint64 `json:"menuId" apitype:"string"`
    Name        string               `json:"name"`
    Description string               `json:"description"`
    Resource    string               `json:"resource"`
    Action      string               `json:"action"`
}

type PermissionUpdateParams struct {
    ID jsontype.SafeUint64 `json:"id" apitype:"string"`
    PermissionCreateParams
}
