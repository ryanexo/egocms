package persistence

import (
    `cms/internal/infra/persistence/gorm/gquery`
    
    `gorm.io/gorm`
)

func NewQuery(db *gorm.DB) *gquery.Query {
    return gquery.Use(db)
}

func NewTransactor(q *gquery.Query) Transactor {
    return q
}
