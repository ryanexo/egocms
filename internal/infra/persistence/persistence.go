package persistence

import (
    `cms/internal/infra/persistence/query`

    `gorm.io/gorm`
)

func NewQuery(db *gorm.DB) *query.Query {
    return query.Use(db)
}

func NewTransactor(q *query.Query) Transactor {
    return q
}
