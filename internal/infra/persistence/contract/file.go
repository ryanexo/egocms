package contract

import (
    `context`
    
    `cms/internal/infra/persistence/datatype`
    `cms/internal/infra/persistence/model`
    
    `gorm.io/gen`
)

type FileRepo interface {
    Repository[FileRepo]
    Create(ctx context.Context, data *model.File) error
    Delete(ctx context.Context, id datatype.SafeUint64) (gen.ResultInfo, error)
    FindByID(ctx context.Context, id datatype.SafeUint64) (*model.File, error)
    FindByPath(ctx context.Context, path string) (*model.File, error)
}
