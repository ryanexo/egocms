package persistence

import (
    `cms/internal/infra/persistence/contract`
    `cms/internal/infra/persistence/query`
    
    `gorm.io/gorm`
)

func NewQuery(db *gorm.DB) *query.Query {
    return query.Use(db)
}

func NewTxManager(q *query.Query) contract.Transactor {
    return q
}
