package api

import (
	"cms/internal/public/apitype"
	"cms/internal/public/jsontype"
)

type CategoryMetaParams struct {
	Title       string   `json:"title"`
	Keywords    []string `json:"keywords"`
	Description string   `json:"description"`
	Thumb       string   `json:"thumb"`
}

type CategoryCreateParams struct {
	ParentID jsontype.SafeUint64 `json:"parentId"`
	Sequence jsontype.SafeInt64  `json:"sequence"`
	Name     string              `json:"name"`
	Path     string              `json:"path"`
	Meta     CategoryMetaParams  `json:"meta"`
}

type CategoryUpdateParams struct {
	ID       jsontype.SafeUint64 `json:"id"`
	Sequence jsontype.SafeInt64  `json:"sequence"`
	Name     string              `json:"name"`
	Path     string              `json:"path"`
	Meta     CategoryMetaParams  `json:"meta"`
}

type CategoryMoveParams struct {
	ID       jsontype.SafeUint64 `json:"id"`
	TargetID jsontype.SafeUint64 `json:"targetId"`
}

type CategoryListParams struct {
	apitype.Pagination

	ID          *jsontype.SafeUint64 `json:"id"`
	ParentID    *jsontype.SafeUint64 `json:"parentId"`
	Name        *string              `json:"name"`
	Title       *string              `json:"title"`
	Keywords    []string             `json:"keywords"`
	Description *string              `json:"description"`
}
