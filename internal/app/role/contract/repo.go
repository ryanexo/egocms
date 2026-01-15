package contract

import (
    `context`
    
    `cms/internal/app/role/internal/dto`
    `cms/internal/infra/persistence/contract`
    `cms/internal/infra/persistence/datatype`
    `cms/internal/infra/persistence/model`
    
    `gorm.io/gen`
)

type RoleRepo interface {
    contract.Repository[RoleRepo]
    Create(ctx context.Context, data *model.Role) error
    Update(ctx context.Context, data *model.Role) (gen.ResultInfo, error)
    FindByID(ctx context.Context, id datatype.SafeUint64) (*model.Role, error)
    Delete(ctx context.Context, id datatype.SafeUint64) (gen.ResultInfo, error)
    List(ctx context.Context, params *dto.RoleListParams) ([]*model.Role, int64, error)
}
