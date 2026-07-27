package contract

import (
    "context"
    
    "cms/internal/infra/persistence"
    `cms/internal/infra/persistence/gorm/model`
    "cms/internal/pkg/datatype"
    
    "cms/internal/modules/role/internal/dto"
    
    "gorm.io/gen"
)

type RoleRepo interface {
    persistence.Repository[RoleRepo]
    Create(ctx context.Context, data *model.Role) error
    Update(ctx context.Context, data *model.Role) (gen.ResultInfo, error)
    FindByID(ctx context.Context, id datatype.SafeUint64) (*model.Role, error)
    Delete(ctx context.Context, id datatype.SafeUint64) (gen.ResultInfo, error)
    List(ctx context.Context, params *dto.RoleListParams) ([]*model.Role, int64, error)
}
