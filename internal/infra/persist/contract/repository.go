package contract

import `cms/internal/infra/persist/query`

type Repository[T any] interface {
    CloneWithQuery(*query.Query) T
}
