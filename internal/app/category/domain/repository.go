package domain

import (
    `context`
    
    `cms/internal/kernel/port`
)

type CategoryRepo interface {
    port.DomainRepository[Category]
    
    CreateSubtree(ctx context.Context, id uint64) error
    Move(ctx context.Context, category *Category, ancestor *Category) error
}
