package contract

import `cms/internal/infra/persistence/query`

type Repository[T any] interface {
    CloneWithQuery(*query.Query) T
}
