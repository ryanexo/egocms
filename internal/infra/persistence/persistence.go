package persistence

import (
    `dpcms/internal/infra/persistence/query`
    
    `gorm.io/gorm`
)

func New(db *gorm.DB) *query.Query {
    return query.Use(db)
}
