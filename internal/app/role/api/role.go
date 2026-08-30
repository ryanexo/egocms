package api

import (
	"cms/internal/public/apitype"
	"cms/internal/public/jsontype"
)

type RoleCreateParams struct {
	Name        string                `json:"name" validate:"required,max=64" label:"名称"`
	Description string                `json:"description" validate:"required,max=255" label:"描述"`
	InheritList []jsontype.SafeUint64 `json:"inheritList" validate:"max=10" label:"继承角色" apitype:"array,string"`
}

type RoleUpdateParams struct {
	apitype.ResourceID
	RoleCreateParams
}

type RoleListParams struct {
	apitype.Pagination
	Name        *string `json:"name"`
	Description *string `json:"description"`
}
