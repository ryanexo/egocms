package dto

import (
	"dpcms/internal/infra/persistence/datatype"
	"dpcms/internal/infra/persistence/dbscope"
	"dpcms/internal/infra/persistence/model"
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
	*model.Role
	InheritList []*model.Role `json:"inheritList"`
}

type RoleListParams struct {
	dbscope.Pagination
	InheritId   *datatype.SafeUint64 `json:"inheritId" swaggertype:"array,string"`
	Name        *string              `json:"name"`
	Description *string              `json:"description"`
}
