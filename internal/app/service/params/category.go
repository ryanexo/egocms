package pcategory

import (
    `gorm.io/gorm`
)

type PCategorySEO struct {
    Title       string `json:"title"`
    Keywords    string `json:"keywords"`
    Description string `json:"description"`
}

type PCreateCategory struct {
    ParentID uint64       `json:"parentID,omitempty"`
    Sequence uint         `json:"sequence,omitempty"`
    Type     uint         `json:"type,omitempty"`
    Name     string       `validate:"required" json:"name" label:"名称"`
    Path     string       `validate:"required" json:"path" label:"路径"`
    Display  uint         `json:"display,omitempty"`
    SEO      PCategorySEO `json:"seo"`
}

type PCategory struct {
    gorm.Model
    PCreateCategory
}
