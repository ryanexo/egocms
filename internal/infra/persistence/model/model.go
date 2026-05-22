package model

import (
    `time`

    `gorm.io/gorm`
)

type Base struct {
	ID        uint64         `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"type:datetime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"type:datetime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
