package model

import (
	"time"

	"gorm.io/gorm"
)

type Base struct {
	ID        uint64         `gorm:"type:bigint unsigned;primaryKey"`
	CreatedAt time.Time      `gorm:"type:datetime"`
	UpdatedAt time.Time      `gorm:"type:datetime"`
	DeletedAt gorm.DeletedAt `gorm:"type:datetime"`
}
