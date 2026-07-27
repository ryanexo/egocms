package dto

import (
	"cms/internal/infra/persistence/datatype"
	datatype2 "cms/internal/pkg/datatype"
	"cms/internal/util/types"
)

type Menu struct {
	types.Base
	ParentID    datatype2.SafeUint64 `json:"parentId"`
	Type        int8                 `json:"type"`
	Name        string               `json:"name"`
	Affix       datatype.BoolInt8    `json:"affix"`
	Icon        string               `json:"icon"`
	ExternalURL string               `json:"externalUrl"`
	Sequence    datatype2.SafeInt64  `json:"sequence"`
	Visible     datatype.BoolInt8    `json:"visible,omitempty" swaggertype:"boolean"`
	URI         string               `json:"uri"`
	Template    string               `json:"template"`
	Remark      string               `json:"remark"`
}

type MenuListQueryParams struct {
	types.Pagination
	Name     *string               `json:"name"`
	ParentID *datatype2.SafeUint64 `json:"parentId" swaggertype:"string"`
}

type MenuCreateParams struct {
	ParentID    datatype2.SafeUint64 `json:"parentId" swaggertype:"string"`
	Type        int8                 `json:"type"`
	Name        string               `validate:"required" json:"name" label:"名称"`
	Affix       datatype.BoolInt8    `json:"affix"`
	Icon        string               `json:"icon"`
	ExternalURL string               `json:"externalUrl"`
	Sequence    datatype2.SafeInt64  `json:"sequence" swaggertype:"string"`
	URI         string               `validate:"required,alphanum" json:"uri"`
	Template    string               `json:"template"`
	Remark      string               `json:"remark"`
}

type MenuUpdateParams struct {
	ID          datatype2.SafeUint64 `validate:"required" json:"id" swaggertype:"string"`
	Type        int8                 `json:"type"`
	Name        string               `json:"name"`
	Affix       datatype.BoolInt8    `json:"affix"`
	Icon        string               `json:"icon"`
	ExternalURL string               `json:"externalUrl"`
	Sequence    datatype2.SafeInt64  `json:"sequence" swaggertype:"string"`
	URI         string               `json:"uri"`
	Template    string               `json:"template"`
	Remark      string               `json:"remark"`
}

type MenuMoveParams struct {
	ID       datatype2.SafeUint64 `validate:"required" json:"id" swaggertype:"string"`
	TargetID datatype2.SafeUint64 `validate:"required" json:"targetId" swaggertype:"string"`
}
