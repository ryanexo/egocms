package persistence

import (
	"cms/internal/infra/persistence/query"
	"database/sql"
)

type Transactor interface {
	Transaction(fc func(tx *query.Query) error, opts ...*sql.TxOptions) error
}

type Repository[T any] interface {
	CloneWithQuery(*query.Query) T
}
