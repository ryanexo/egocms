package domain

import (
    `context`
)

type CategoryRepo interface {
    Create(ctx context.Context, category Category) error
    CreateSubtree(ctx context.Context, category Category) error
    Update(ctx context.Context, category Category) error
    Move(ctx context.Context, category Category, ancestor Category) error
    
    Delete(ctx context.Context, category Category) error
    Find(ctx context.Context, id uint64) (Category, error)
}
