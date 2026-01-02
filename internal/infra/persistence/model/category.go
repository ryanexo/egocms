package model

import `dpcms/internal/infra/persistence/datatype`

type CategoryContext struct {
    ID datatype.SafeUint64 `gorm:"primaryKey"`
    ClosureTableModel
}

type Category struct {
    Base
    ParentID datatype.SafeUint64 `gorm:"not null;index" json:"parentId"`
    Sequence uint                `gorm:"not null;index" json:"sequence"`
    Name     string              `gorm:"type:varchar(255);not null" json:"name"`
    Path     string              `gorm:"type:varchar(64);not null;uniqueIndex" json:"path"`
    Type     uint                `gorm:"type:tinyint;not null;comment:'0:普通分类,1:单页型分类,2:链接'" json:"type"`
    Display  uint                `gorm:"type:tinyint;not null" json:"display"`
    SEO      *CategorySeo        `gorm:"foreignKey:CategoryID;references:ID" json:"SEO,omitempty"`
}

type CategorySeo struct {
    Base
    CategoryID  datatype.SafeUint64 `gorm:"not null;index" json:"categoryId"`
    Title       string              `gorm:"type:varchar(255);not null" json:"title"`
    Keywords    string              `gorm:"type:varchar(255);not null" json:"keywords"`
    Description string              `gorm:"type:varchar(255);not null" json:"description"`
}
