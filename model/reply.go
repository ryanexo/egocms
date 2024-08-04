package model

import (
	"gorm.io/gorm"
)

type Reply struct {
	gorm.Model
	TargetID uint64 `gorm:"not null"`
	UserID   uint64 `gorm:"not null"`
	Content  string `gorm:"text;not null"`
}
