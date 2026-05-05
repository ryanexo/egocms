package port

import `context`

type DomainRepository[T any] interface {
    Create(ctx context.Context, category *T) (uint64, error)
    Update(ctx context.Context, category *T) error
    Delete(ctx context.Context, category *T) error
    Find(ctx context.Context, id uint64) (*T, error)
}
