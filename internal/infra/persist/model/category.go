package model

import `cms/internal/infra/persist/datatype`

type CategoryContext struct {
    ID datatype.SafeUint64 `gorm:"primaryKey"`
    ClosureTableModel
}

type Category struct {
    Base
    ParentID *datatype.SafeUint64 `gorm:"index"`
    Sequence datatype.SafeInt64   `gorm:"not null;index"`
    Name     string               `gorm:"type:varchar(255);not null"`
    Path     string               `gorm:"type:varchar(64);not null;uniqueIndex"`
    Type     int8                 `gorm:"type:tinyint;not null;comment:'0:普通分类,1:单页型分类,2:链接'"`
    URL      *string              `gorm:"type:varchar(255);default:null"`
    Visible  datatype.BoolInt8    `gorm:"type:tinyint;not null"`
    SEO      *CategorySeo         `gorm:"foreignKey:CategoryID;references:ID"`
}

type CategorySeo struct {
    Base
    CategoryID  datatype.SafeUint64 `gorm:"not null;index"`
    Title       string              `gorm:"type:varchar(255);not null"`
    Keywords    string              `gorm:"type:varchar(255);not null"`
    Description string              `gorm:"type:varchar(255);not null"`
}

type CategoryPageContent struct {
    Base
    CategoryID datatype.SafeUint64 `gorm:"not null;index"`
    Content    *string             `gorm:"type:text;default:null"`
}
