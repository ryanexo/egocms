package contract

import (
    "context"
    
    "cms/internal/infra/persistence"
    `cms/internal/infra/persistence/gorm/model`
    "cms/internal/pkg/datatype"
    
    "gorm.io/gen"
)

type PermissionRepo interface {
    persistence.Repository[PermissionRepo]
    Create(ctx context.Context, data *model.Permission) error
    Update(ctx context.Context, data *model.Permission) (gen.ResultInfo, error)
    Delete(ctx context.Context, id datatype.SafeUint64) (gen.ResultInfo, error)
    FindByID(ctx context.Context, id datatype.SafeUint64) (*model.Permission, error)
    FindByMenuID(ctx context.Context, id datatype.SafeUint64) ([]*model.Permission, error)
    FindByResource(ctx context.Context, resource string) ([]*model.Permission, error)
    FindNoMenuID(ctx context.Context) ([]*model.Permission, error)
    FindAll(ctx context.Context) ([]*model.Permission, error)
}
