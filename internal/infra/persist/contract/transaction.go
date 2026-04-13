package contract

import (
    `database/sql`
    
    `cms/internal/infra/persist/query`
)

type Transactor interface {
    Transaction(fc func(tx *query.Query) error, opts ...*sql.TxOptions) error
}
