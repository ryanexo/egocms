package service

import (
    `context`
    
    `dpcms/internal/infra/persistence/contract`
    `dpcms/internal/infra/persistence/datatype`
    `dpcms/internal/infra/persistence/model`
    
    `gorm.io/gen`
)

type FileRepo interface {
    contract.Repository[FileRepo]
    Create(ctx context.Context, data *model.File) error
    Delete(ctx context.Context, id datatype.SafeUint64) (gen.ResultInfo, error)
    FindByID(ctx context.Context, id datatype.SafeUint64) (*model.File, error)
    FindByPath(ctx context.Context, path string) (*model.File, error)
}
