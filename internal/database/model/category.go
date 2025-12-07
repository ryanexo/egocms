package model

import (
    `dpcms/internal/packages/database`
)

type CategoryContext struct {
    ID int64 `gorm:"primaryKey"`
    ClosureTableModel
}

type Category struct {
    database.Model
    ParentID int64        `gorm:"not null;index"`
    Sequence uint         `gorm:"not null;index"`
    Name     string       `gorm:"type:varchar(255);not null"`
    Path     string       `gorm:"type:varchar(64);not null;uniqueIndex"`
    Type     uint         `gorm:"type:tinyint;not null;comment:'0:普通分类,1:单页型分类,2:链接'"`
    Display  uint         `gorm:"type:tinyint;not null"`
    SEO      *CategorySeo `gorm:"foreignKey:CategoryID;references:ID"`
    Children []Category   `gorm:"foreignKey:ParentID;references:ID"`
}

type CategorySeo struct {
    database.Model
    CategoryID  int64  `gorm:"not null;index"`
    Title       string `gorm:"type:varchar(255);not null"`
    Keywords    string `gorm:"type:varchar(255);not null"`
    Description string `gorm:"type:varchar(255);not null"`
}
