package contract

import `dpcms/internal/infra/persistence/query`

type Repository[T any] interface {
    CloneWithQuery(*query.Query) T
}
