package service

import (
    `context`
    
    `dpcms/internal/app/role/internal/dto`
    `dpcms/internal/infra/persistence/contract`
    `dpcms/internal/infra/persistence/datatype`
    `dpcms/internal/infra/persistence/model`
    
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
