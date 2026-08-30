package contract

import (
    "context"
    
    `cms/internal/public/jsontype`
    `cms/internal/public/model`
    
    "gorm.io/gen"
)

type FileRepo interface {
    persistence.Repository[FileRepo]
    Create(ctx context.Context, data *model.File) error
    Delete(ctx context.Context, id jsontype.SafeUint64) (gen.ResultInfo, error)
    FindByID(ctx context.Context, id jsontype.SafeUint64) (*model.File, error)
    FindByPath(ctx context.Context, path string) (*model.File, error)
}
