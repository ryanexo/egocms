package entity

import (
    `gorm.io/gorm`
)

type Category struct {
    gorm.Model
    ID             uint64 `gorm:"primaryKey"`
    Sequence       uint64 `gorm:"not null"`
    RootID         uint64 `gorm:"not null;index"`
    ParentID       uint64 `gorm:"not null;index"`
    Name           string `gorm:"type:varchar(255);not null"`
    Alias          string `gorm:"type:varchar(64);not null;uniqueIndex"`
    SeoTitle       string `gorm:"type:varchar(255);not null"`
    SeoKeywords    string `gorm:"type:varchar(255);not null"`
    SeoDescription string `gorm:"type:varchar(255);not null"`
    Type           uint8  `gorm:"not null"`
    Display        uint8  `gorm:"not null"`
}
