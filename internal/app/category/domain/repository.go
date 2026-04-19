package domain

import (
    `cms/internal/infra/persist/contract`
)

type CategoryRepo interface {
    contract.Repository[CategoryRepo]
}
