package model

import (
    `time`
    
    `gorm.io/gorm`
)

type Base struct {
    ID        uint64
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt
}
