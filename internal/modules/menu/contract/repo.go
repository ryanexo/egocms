package contract

import (
    "context"
    
    `cms/internal/public/jsontype`
    `cms/internal/public/model`
    
    "cms/internal/modules/menu/internal/dto"
    
    "gorm.io/gen"
)

type MenuRepo interface {
    persistence.Repository[MenuRepo]
    Create(ctx context.Context, menu *model.Menu) error
    CreateSubtree(ctx context.Context, id jsontype.SafeUint64, parentID jsontype.SafeUint64) error
    Update(ctx context.Context, data *model.Menu) (gen.ResultInfo, error)
    Move(ctx context.Context, fromNode jsontype.SafeUint64, toNode jsontype.SafeUint64) error
    Delete(ctx context.Context, id jsontype.SafeUint64) error
    FindByID(ctx context.Context, id jsontype.SafeUint64) (*model.Menu, error)
    FindByIDWithAncestor(ctx context.Context, ancestor jsontype.SafeUint64, descendant jsontype.SafeUint64) (*model.MenuContext, error)
    List(ctx context.Context, params dto.MenuListQueryParams) ([]*model.Menu, int64, error)
}
