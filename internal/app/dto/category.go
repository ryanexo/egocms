package dto

import (
    dtotype `dpcms/internal/app/dto/type`
    `dpcms/internal/infra/persistence/datatype`
    `dpcms/internal/infra/persistence/dbscope`
)

type Category struct {
    dtotype.Base `json:"dtotype.Base"`
    ParentID     datatype.SafeUint64 `json:"parentId"`
    Sequence     uint                `json:"sequence"`
    Name         string              `json:"name"`
    Path         string              `json:"path"`
    Type         uint                `json:"type"`
    Display      uint                `json:"display"`
    SEO          *CategorySEO        `json:"seo"`
}

type CategorySEO struct {
    Title       string `json:"title"`
    Keywords    string `json:"keywords"`
    Description string `json:"description"`
}

type CategoryCreateParams struct {
    ParentID datatype.SafeUint64 `json:"parentID,omitempty" swaggertype:"string"`
    Sequence uint                `json:"sequence,omitempty"`
    Type     uint                `json:"type,omitempty"`
    Name     string              `validate:"required" json:"name" label:"名称"`
    Path     string              `validate:"required" json:"path" label:"路径"`
    Display  uint                `json:"display,omitempty"`
    SEO      CategorySEO         `json:"seo"`
}

type CategoryUpdateParams struct {
    ID       datatype.SafeUint64 `validate:"required" json:"id" swaggertype:"string"`
    Sequence uint                `json:"sequence,omitempty"`
    Type     uint                `json:"type,omitempty"`
    Name     string              `json:"name"`
    Path     string              `json:"path"`
    Display  uint                `json:"display,omitempty"`
    SEO      CategorySEO         `json:"seo"`
}

type CategoryMoveParams struct {
    ID       datatype.SafeUint64 `json:"id" swaggertype:"string"`
    TargetID datatype.SafeUint64 `json:"targetId" swaggertype:"string"`
}

type CategoryListParams struct {
    dbscope.Pagination
    ID       *datatype.SafeUint64 `json:"id" swaggertype:"string"`
    ParentID *datatype.SafeUint64 `json:"parentId" swaggertype:"string"`
    Type     *uint                `json:"type"`
    Name     *string              `json:"name"`
    Path     *string              `json:"path"`
    Display  *uint                `json:"display"`
}
