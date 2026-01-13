package adapter

import (
    `context`
    
    `dpcms/internal/app/setting/service`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/infra/persistence/query`
    
    `gorm.io/gen`
)

type settingRepo struct {
    persist *query.Query
}

func NewConfigRepo(persist *query.Query) service.SettingRepo {
    return &settingRepo{persist}
}

func (s *settingRepo) CloneWithQuery(q *query.Query) service.SettingRepo {
    return NewConfigRepo(q)
}

func (s *settingRepo) Add(ctx context.Context, data *model.Config) error {
    return s.persist.Config.WithContext(ctx).Create(data)
}

func (s *settingRepo) AddInBatches(ctx context.Context, data []*model.Config) error {
    return s.persist.Config.WithContext(ctx).CreateInBatches(data, 500)
}

func (s *settingRepo) Update(ctx context.Context, data *model.Config) (gen.ResultInfo, error) {
    dao := s.persist.Config
    return dao.WithContext(ctx).Where(dao.Key.Eq(data.Key)).Update(dao.Value, data.Value)
}

func (s *settingRepo) Remove(ctx context.Context, key string) (gen.ResultInfo, error) {
    dao := s.persist.Config
    return dao.WithContext(ctx).Where(dao.Key.Eq(key)).Delete()
}

func (s *settingRepo) Get(ctx context.Context, key string) (string, error) {
    dao := s.persist.Config
    data, err := dao.WithContext(ctx).Where(dao.Key.Eq(key)).First()
    if err != nil {
        return "", err
    }
    return data.Value, nil
}

func (s *settingRepo) GetAll(ctx context.Context) ([]*model.Config, error) {
    dao := s.persist.Config
    return dao.WithContext(ctx).Find()
}
