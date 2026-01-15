package persist

import (
    `cms/internal/infra/persist/contract`
    `cms/internal/infra/persist/query`
    
    `gorm.io/gorm`
)

func NewQuery(db *gorm.DB) *query.Query {
    return query.Use(db)
}

func NewTxManager(q *query.Query) contract.TxManager {
    return q
}
