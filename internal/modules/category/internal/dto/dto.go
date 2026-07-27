package dto

import (
	"cms/internal/infra/persistence/datatype"
	datatype2 "cms/internal/pkg/datatype"
	"cms/internal/util/types"
)

type Category struct {
	types.Base
	ParentID datatype2.SafeUint64 `json:"parentId"`
	Sequence datatype2.SafeInt64  `json:"sequence"`
	Name     string               `json:"name"`
	Path     string               `json:"path"`
	Type     int8                 `json:"type"`
	Visible  *datatype.BoolInt8   `json:"visible,omitempty" swaggertype:"boolean"`
	SEO      *CategorySEO         `json:"seo"`
}

type CategorySEO struct {
	Title       string `json:"title"`
	Keywords    string `json:"keywords"`
	Description string `json:"description"`
}

type CategoryCreateParams struct {
	ParentID datatype2.SafeUint64 `json:"parentID" swaggertype:"string"`
	Sequence datatype2.SafeInt64  `json:"sequence"`
	Type     int8                 `json:"type"`
	Name     string               `validate:"required" json:"name" label:"名称"`
	Path     string               `validate:"required" json:"path" label:"路径"`
	Visible  *datatype.BoolInt8   `json:"visible" swaggertype:"boolean"`
	SEO      CategorySEO          `json:"seo"`
}

type CategoryUpdateParams struct {
	ID       datatype2.SafeUint64 `validate:"required" json:"id" swaggertype:"string"`
	Sequence datatype2.SafeInt64  `json:"sequence"`
	Type     int8                 `json:"type"`
	Name     string               `json:"name"`
	Path     string               `json:"path"`
	Visible  *datatype.BoolInt8   `json:"visible" swaggertype:"boolean"`
	SEO      CategorySEO          `json:"seo"`
}

type CategoryMoveParams struct {
	ID       datatype2.SafeUint64 `json:"id" swaggertype:"string"`
	TargetID datatype2.SafeUint64 `json:"targetId" swaggertype:"string"`
}

type CategoryListParams struct {
	types.Pagination
	ID       *datatype2.SafeUint64 `json:"id" swaggertype:"string"`
	ParentID *datatype2.SafeUint64 `json:"parentId" swaggertype:"string"`
	Type     *int8                 `json:"type"`
	Name     *string               `json:"name"`
	Path     *string               `json:"path"`
	Visible  *datatype.BoolInt8    `json:"visible" swaggertype:"boolean"`
}
