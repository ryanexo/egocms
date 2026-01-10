package adapter

import (
    `context`
    
    `dpcms/internal/app/config/service`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/infra/persistence/query`
    
    `gorm.io/gen`
)

type configRepo struct {
    persist *query.Query
}

func NewConfigRepo(persist *query.Query) service.ConfigRepo {
    return &configRepo{persist}
}

func (s *configRepo) CloneWithQuery(q *query.Query) service.ConfigRepo {
    return NewConfigRepo(q)
}

func (s *configRepo) Add(ctx context.Context, data *model.Config) error {
    return s.persist.Config.WithContext(ctx).Create(data)
}

func (s *configRepo) AddInBatches(ctx context.Context, data []*model.Config) error {
    return s.persist.Config.WithContext(ctx).CreateInBatches(data, 500)
}

func (s *configRepo) Update(ctx context.Context, data *model.Config) (gen.ResultInfo, error) {
    dao := s.persist.Config
    return dao.WithContext(ctx).Where(dao.Key.Eq(data.Key)).Update(dao.Value, data.Value)
}

func (s *configRepo) Remove(ctx context.Context, key string) (gen.ResultInfo, error) {
    dao := s.persist.Config
    return dao.WithContext(ctx).Where(dao.Key.Eq(key)).Delete()
}

func (s *configRepo) Get(ctx context.Context, key string) (string, error) {
    dao := s.persist.Config
    data, err := dao.WithContext(ctx).Where(dao.Key.Eq(key)).First()
    if err != nil {
        return "", err
    }
    return data.Value, nil
}

func (s *configRepo) GetAll(ctx context.Context) ([]*model.Config, error) {
    dao := s.persist.Config
    return dao.WithContext(ctx).Find()
}
