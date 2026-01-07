package dto

import (
    `dpcms/internal/infra/persistence/datatype`
    `dpcms/internal/types`
)

type Permission struct {
    types.Base
    MenuID      *datatype.SafeUint64 `json:"menuId" swaggertype:"string"`
    Name        string               `json:"name"`
    Description string               `json:"description"`
    Resource    string               `json:"resource"`
    Action      string               `json:"action"`
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
