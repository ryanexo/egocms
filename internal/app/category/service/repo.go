package service

import (
    `context`
    
    `dpcms/internal/app/category/internal/dto`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/infra/persistence/contract`
    
    `gorm.io/gen`
)

type CategoryRepo interface {
    contract.Repository[CategoryRepo]
    Create(ctx context.Context, category *model.Category) error
    CreateSubtree(ctx context.Context, id uint64, parentID uint64) error
    Update(ctx context.Context, data *model.Category) (gen.ResultInfo, error)
    Move(ctx context.Context, fromNode uint64, toNode uint64) error
    Delete(ctx context.Context, id uint64) error
    FindByID(ctx context.Context, id uint64) (*model.Category, error)
    FindByIDWithAncestor(ctx context.Context, ancestor uint64, descendant uint64) (*model.CategoryContext, error)
    List(ctx context.Context, params dto.CategoryListParams) ([]*model.Category, int64, error)
}
