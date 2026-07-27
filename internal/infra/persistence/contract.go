package persistence

import (
    "database/sql"
    
    `cms/internal/infra/persistence/gorm/gquery`
)

type Transactor interface {
    Transaction(fc func(tx *gquery.Query) error, opts ...*sql.TxOptions) error
}

type Repository[T any] interface {
    CloneWithQuery(*gquery.Query) T
}
