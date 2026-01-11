package service

import (
    `context`
    
    `dpcms/internal/app/menu/internal/dto`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/infra/persistence/contract`
    
    `gorm.io/gen`
)

type MenuRepo interface {
    contract.Repository[MenuRepo]
    Create(ctx context.Context, menu *model.Menu) error
    CreateSubtree(ctx context.Context, id uint64, parentID uint64) error
    Update(ctx context.Context, data *model.Menu) (gen.ResultInfo, error)
    Move(ctx context.Context, fromNode uint64, toNode uint64) error
    Delete(ctx context.Context, id uint64) error
    FindByID(ctx context.Context, id uint64) (*model.Menu, error)
    FindByIDWithAncestor(ctx context.Context, ancestor uint64, descendant uint64) (*model.MenuContext, error)
    List(ctx context.Context, params dto.MenuListQueryParams) ([]*model.Menu, int64, error)
}
