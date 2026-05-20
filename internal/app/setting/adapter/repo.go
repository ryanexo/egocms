package adapter

import (
    `context`
    
    `cms/internal/app/setting/contract`
    `cms/internal/infra/persistence/model`
    `cms/internal/infra/persistence/query`
    
    `gorm.io/gen`
)

type settingRepo struct {
    persist *query.Query
}

func NewConfigRepo(persist *query.Query) contract.SettingRepo {
    return &settingRepo{persist: persist}
}

func (s *settingRepo) CloneWithQuery(q *query.Query) contract.SettingRepo {
    return NewConfigRepo(q)
}

func (s *settingRepo) Add(ctx context.Context, data *model.Setting) error {
    return s.persist.Setting.WithContext(ctx).Create(data)
}

func (s *settingRepo) AddInBatches(ctx context.Context, data []*model.Setting) error {
    return s.persist.Setting.WithContext(ctx).CreateInBatches(data, 500)
}

func (s *settingRepo) Update(ctx context.Context, data *model.Setting) (gen.ResultInfo, error) {
    dao := s.persist.Setting
    return dao.WithContext(ctx).Where(dao.Key.Eq(data.Key)).Update(dao.Value, data.Value)
}

func (s *settingRepo) Remove(ctx context.Context, key string) (gen.ResultInfo, error) {
    dao := s.persist.Setting
    return dao.WithContext(ctx).Where(dao.Key.Eq(key)).Delete()
}

func (s *settingRepo) Get(ctx context.Context, key string) (string, error) {
    dao := s.persist.Setting
    data, err := dao.WithContext(ctx).Where(dao.Key.Eq(key)).First()
    if err != nil {
        return "", err
    }
    return data.Value, nil
}

func (s *settingRepo) GetAll(ctx context.Context) ([]*model.Setting, error) {
    return s.persist.Setting.WithContext(ctx).Find()
}
