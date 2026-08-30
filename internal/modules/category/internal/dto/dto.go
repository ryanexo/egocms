package dto

import (
    "cms/internal/infra/store/datatype"
    `cms/internal/public/apitype`
    `cms/internal/public/jsontype`
)

type Category struct {
    apitype.Base
    ParentID jsontype.SafeUint64 `json:"parentId"`
    Sequence jsontype.SafeInt64  `json:"sequence"`
    Name     string              `json:"name"`
    Path     string              `json:"path"`
    Type     int8                `json:"type"`
    Visible  *datatype.BoolInt8  `json:"visible,omitempty" apitype:"boolean"`
    SEO      *CategorySEO        `json:"seo"`
}

type CategorySEO struct {
    Title       string `json:"title"`
    Keywords    string `json:"keywords"`
    Description string `json:"description"`
}

type CategoryCreateParams struct {
    ParentID jsontype.SafeUint64 `json:"parentID" apitype:"string"`
    Sequence jsontype.SafeInt64  `json:"sequence"`
    Type     int8                `json:"type"`
    Name     string              `validate:"required" json:"name" label:"名称"`
    Path     string              `validate:"required" json:"path" label:"路径"`
    Visible  *datatype.BoolInt8  `json:"visible" apitype:"boolean"`
    SEO      CategorySEO         `json:"seo"`
}

type CategoryUpdateParams struct {
    ID       jsontype.SafeUint64 `validate:"required" json:"id" apitype:"string"`
    Sequence jsontype.SafeInt64  `json:"sequence"`
    Type     int8                `json:"type"`
    Name     string              `json:"name"`
    Path     string              `json:"path"`
    Visible  *datatype.BoolInt8  `json:"visible" apitype:"boolean"`
    SEO      CategorySEO         `json:"seo"`
}

type CategoryMoveParams struct {
    ID       jsontype.SafeUint64 `json:"id" apitype:"string"`
    TargetID jsontype.SafeUint64 `json:"targetId" apitype:"string"`
}

type CategoryListParams struct {
    apitype.Pagination
    ID       *jsontype.SafeUint64 `json:"id" apitype:"string"`
    ParentID *jsontype.SafeUint64 `json:"parentId" apitype:"string"`
    Type     *int8                `json:"type"`
    Name     *string              `json:"name"`
    Path     *string              `json:"path"`
    Visible  *datatype.BoolInt8   `json:"visible" apitype:"boolean"`
}
