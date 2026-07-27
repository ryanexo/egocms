package contract

import (
    "context"
    
    "cms/internal/infra/persistence"
    `cms/internal/infra/persistence/gorm/model`
    "cms/internal/pkg/datatype"
    
    "gorm.io/gen"
)

type FileRepo interface {
    persistence.Repository[FileRepo]
    Create(ctx context.Context, data *model.File) error
    Delete(ctx context.Context, id datatype.SafeUint64) (gen.ResultInfo, error)
    FindByID(ctx context.Context, id datatype.SafeUint64) (*model.File, error)
    FindByPath(ctx context.Context, path string) (*model.File, error)
}
