package contract

import (
    `cms/internal/app/category/internal`
    
    "context"
    
    "cms/internal/modules/category/internal/dto"
    
    "gorm.io/gen"
)

type CategoryRepo interface {
    persistence.Repository[CategoryRepo]
    Create(ctx context.Context, category *internal.Category) error
    CreateSubtree(ctx context.Context, id uint64, parentID uint64) error
    Update(ctx context.Context, data *internal.Category) (gen.ResultInfo, error)
    Move(ctx context.Context, fromNode uint64, toNode uint64) error
    Delete(ctx context.Context, id uint64) error
    FindByID(ctx context.Context, id uint64) (*internal.Category, error)
    FindByIDWithAncestor(ctx context.Context, ancestor uint64, descendant uint64) (*internal.CategoryContext, error)
    List(ctx context.Context, params dto.CategoryListParams) ([]*internal.Category, int64, error)
}
