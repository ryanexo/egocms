package model

import (
    `time`
    
    `cms/internal/infra/persistence/datatype`
    
    `gorm.io/gorm`
)

type Base struct {
    ID        datatype.SafeUint64 `gorm:"primaryKey" json:"id"`
    CreatedAt time.Time           `gorm:"type:datetime" json:"createdAt"`
    UpdatedAt time.Time           `gorm:"type:datetime" json:"updatedAt"`
    DeletedAt gorm.DeletedAt      `gorm:"index" json:"-"`
}
