package contract

import (
	"cms/internal/infra/persistence"
	"cms/internal/pkg/datatype"
	"context"

	"cms/internal/app/menu/internal/dto"
	"cms/internal/infra/persistence/model"

	"gorm.io/gen"
)

type MenuRepo interface {
	persistence.Repository[MenuRepo]
	Create(ctx context.Context, menu *model.Menu) error
	CreateSubtree(ctx context.Context, id datatype.SafeUint64, parentID datatype.SafeUint64) error
	Update(ctx context.Context, data *model.Menu) (gen.ResultInfo, error)
	Move(ctx context.Context, fromNode datatype.SafeUint64, toNode datatype.SafeUint64) error
	Delete(ctx context.Context, id datatype.SafeUint64) error
	FindByID(ctx context.Context, id datatype.SafeUint64) (*model.Menu, error)
	FindByIDWithAncestor(ctx context.Context, ancestor datatype.SafeUint64, descendant datatype.SafeUint64) (*model.MenuContext, error)
	List(ctx context.Context, params dto.MenuListQueryParams) ([]*model.Menu, int64, error)
}
