package service

import (
    `context`
    
    `dpcms/internal/infra/persistence/contract`
    `dpcms/internal/infra/persistence/model`
    
    `gorm.io/gen`
)

type ConfigRepo interface {
    contract.Repository[ConfigRepo]
    Add(ctx context.Context, data *model.Config) error
    AddInBatches(ctx context.Context, data []*model.Config) error
    Update(ctx context.Context, data *model.Config) (gen.ResultInfo, error)
    Remove(ctx context.Context, key string) (gen.ResultInfo, error)
    Get(ctx context.Context, key string) (string, error)
    GetAll(ctx context.Context) ([]*model.Config, error)
}
