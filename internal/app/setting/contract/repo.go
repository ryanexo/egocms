package contract

import (
    `context`
    
    `cms/internal/infra/persist/contract`
    `cms/internal/infra/persist/model`
    
    `gorm.io/gen`
)

type SettingRepo interface {
    contract.Repository[SettingRepo]
    Add(ctx context.Context, data *model.Setting) error
    AddInBatches(ctx context.Context, data []*model.Setting) error
    Update(ctx context.Context, data *model.Setting) (gen.ResultInfo, error)
    Remove(ctx context.Context, key string) (gen.ResultInfo, error)
    Get(ctx context.Context, key string) (string, error)
    GetAll(ctx context.Context) ([]*model.Setting, error)
}
