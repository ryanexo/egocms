package persistence

import (
    `dpcms/internal/infra/persistence/contract`
    `dpcms/internal/infra/persistence/query`
    
    `gorm.io/gorm`
)

func NewQuery(db *gorm.DB) *query.Query {
    return query.Use(db)
}

func NewTxManager(q *query.Query) contract.TxManager {
    return q
}
