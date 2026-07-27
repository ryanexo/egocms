package contract

import (
    `cms/internal/infra/persistence/gorm/model`
    
    "context"
    
    "gorm.io/gen"
)

type SettingRepo interface {
    Add(ctx context.Context, data *model.Setting) error
    AddInBatches(ctx context.Context, data []*model.Setting) error
    Update(ctx context.Context, data *model.Setting) (gen.ResultInfo, error)
    Remove(ctx context.Context, key string) (gen.ResultInfo, error)
    Get(ctx context.Context, key string) (string, error)
    GetAll(ctx context.Context) ([]*model.Setting, error)
}
