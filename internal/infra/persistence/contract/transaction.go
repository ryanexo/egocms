package contract

import (
    `database/sql`
    
    `cms/internal/infra/persistence/query`
)

type TxManager interface {
    Transaction(fc func(tx *query.Query) error, opts ...*sql.TxOptions) error
}
