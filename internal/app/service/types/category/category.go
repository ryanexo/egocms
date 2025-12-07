package category

import (
    `dpcms/internal/app/service/types`
)

type SEODetail struct {
    Title       string
    Keywords    string
    Description string
}

type CreateParams struct {
    ParentID int64     `json:"parentID,omitempty"`
    Sequence uint      `json:"sequence,omitempty"`
    Type     uint      `json:"type,omitempty"`
    Name     string    `validate:"required" json:"name" label:"名称"`
    Path     string    `validate:"required" json:"path" label:"路径"`
    Display  uint      `json:"display,omitempty"`
    SEO      SEODetail `json:"seo"`
}

type Detail struct {
    types.Meta
    CreateParams
}
