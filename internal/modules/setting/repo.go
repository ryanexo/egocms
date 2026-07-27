package setting

import (
    `context`
    
    `cms/internal/infra/persistence/gorm/model`
    `cms/internal/infra/persistence/gorm/gquery`
    `cms/internal/modules/setting/contract`
    
    `gorm.io/gen`
)

type configRepo struct {
    query *gquery.Query
}

func NewConfigRepo(persist *gquery.Query) contract.ConfigRepo {
    return &configRepo{query: persist}
}

func (s *configRepo) Add(ctx context.Context, data *model.Config) error {
    return s.query.Config.WithContext(ctx).Create(data)
}

func (s *configRepo) AddInBatches(ctx context.Context, data []*model.Config) error {
    return s.query.Config.WithContext(ctx).CreateInBatches(data, 500)
}

func (s *configRepo) Update(ctx context.Context, data *model.Config) (gen.ResultInfo, error) {
    dao := s.query.Config
    return dao.WithContext(ctx).Where(dao.Key.Eq(data.Key)).Update(dao.Value, data.Value)
}

func (s *configRepo) Remove(ctx context.Context, key string) (gen.ResultInfo, error) {
    dao := s.query.Config
    return dao.WithContext(ctx).Where(dao.Key.Eq(key)).Delete()
}

func (s *configRepo) Get(ctx context.Context, key string) (string, error) {
    dao := s.query.Config
    data, err := dao.WithContext(ctx).Where(dao.Key.Eq(key)).First()
    if err != nil {
        return "", err
    }
    return data.Value, nil
}

func (s *configRepo) GetAll(ctx context.Context) ([]*model.Config, error) {
    return s.query.Config.WithContext(ctx).Find()
}
