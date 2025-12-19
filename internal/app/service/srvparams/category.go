package srvparams

import `dpcms/internal/app/helper/dbscope`

type CategorySEO struct {
    Title       string `json:"title"`
    Keywords    string `json:"keywords"`
    Description string `json:"description"`
}

type CategoryCreateParams struct {
    ParentID uint64      `json:"parentID,omitempty"`
    Sequence uint        `json:"sequence,omitempty"`
    Type     uint        `json:"type,omitempty"`
    Name     string      `validate:"required" json:"name" label:"名称"`
    Path     string      `validate:"required" json:"path" label:"路径"`
    Display  uint        `json:"display,omitempty"`
    SEO      CategorySEO `json:"seo"`
}

type CategoryUpdateParams struct {
    ID       uint64      `validate:"required" json:"id"`
    Sequence uint        `json:"sequence,omitempty"`
    Type     uint        `json:"type,omitempty"`
    Name     string      `json:"name"`
    Path     string      `json:"path"`
    Display  uint        `json:"display,omitempty"`
    SEO      CategorySEO `json:"seo"`
}

type CategoryMoveParams struct {
    ID       uint64 `json:"id"`
    TargetID uint64 `json:"targetId"`
}

type CategoryListParams struct {
    dbscope.Pagination
    ID       *uint64 `json:"id"`
    ParentID *uint64 `json:"parentId"`
    Type     *uint   `json:"type"`
    Name     *string `json:"name"`
    Path     *string `json:"path"`
    Display  *uint   `json:"display"`
}
