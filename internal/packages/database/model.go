package database

import (
    `time`
    
    `gorm.io/gorm`
)

type Model struct {
    ID        uint64         `gorm:"primaryKey" json:"id"`
    CreatedAt time.Time      `json:"createdAt"`
    UpdatedAt time.Time      `json:"updatedAt"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
