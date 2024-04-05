package entity

import (
    `gorm.io/gorm`
)

type Article struct {
    gorm.Model
    User        User
    Category    Category
    ArticleData ArticleData
    UserID      uint64 `gorm:"not null"`
    CategoryID  uint64 `gorm:"not null"`
    Click       uint64 `gorm:"not null"`
    CustomURL   string `gorm:"not null"`
    Title       string `gorm:"type:varchar(255);not null"`
    Description string `gorm:"type:varchar(255);not null"`
    Thumb       string `gorm:"type:varchar(255);not null"`
    SeoKeywords string `gorm:"type:varchar(255);not null"`
    Flag        uint8  `gorm:"not null"`
    Status      uint8  `gorm:"not null"`
}
