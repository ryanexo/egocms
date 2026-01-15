package adapter

import (
    `context`
    
    `cms/internal/app/setting/contract`
    `cms/internal/infra/cache`
    `cms/internal/infra/persistence/model`
    `cms/internal/infra/persistence/query`
    
    `gorm.io/gen`
)

type settingRepo struct {
    persist *query.Query
    cache   *cache.SettingCache
}

func NewConfigRepo(persist *query.Query, c *cache.SettingCache) contract.SettingRepo {
    return &settingRepo{persist: persist, cache: c}
}

func (s *settingRepo) CloneWithQuery(q *query.Query) contract.SettingRepo {
    return NewConfigRepo(q, s.cache)
}

func (s *settingRepo) Add(ctx context.Context, data *model.Setting) error {
    err := s.persist.Setting.WithContext(ctx).Create(data)
    if err == nil {
        s.cache.Set(data.Key, data.Value, 0)
    }
    return err
}

func (s *settingRepo) AddInBatches(ctx context.Context, data []*model.Setting) error {
    err := s.persist.Setting.WithContext(ctx).CreateInBatches(data, 500)
    if err == nil {
        for _, item := range data {
            s.cache.Set(item.Key, item.Value, 0)
        }
    }
    return err
}

func (s *settingRepo) Update(ctx context.Context, data *model.Setting) (gen.ResultInfo, error) {
    dao := s.persist.Setting
    result, err := dao.WithContext(ctx).Where(dao.Key.Eq(data.Key)).Update(dao.Value, data.Value)
    if err == nil {
        s.cache.Set(data.Key, data.Value, 0)
    }
    return result, err
}

func (s *settingRepo) Remove(ctx context.Context, key string) (gen.ResultInfo, error) {
    dao := s.persist.Setting
    result, err := dao.WithContext(ctx).Where(dao.Key.Eq(key)).Delete()
    if err == nil {
        s.cache.Del(key)
    }
    return result, err
}

func (s *settingRepo) Get(ctx context.Context, key string) (string, error) {
    value, ok := s.cache.Get(key)
    if ok {
        return value, nil
    }
    dao := s.persist.Setting
    data, err := dao.WithContext(ctx).Where(dao.Key.Eq(key)).First()
    if err != nil {
        return "", err
    }
    s.cache.Set(data.Key, data.Value, 0)
    return data.Value, nil
}

func (s *settingRepo) GetAll(ctx context.Context) ([]*model.Setting, error) {
    result, err := s.persist.Setting.WithContext(ctx).Find()
    if err != nil {
        return nil, err
    }
    for _, item := range result {
        s.cache.Set(item.Key, item.Value, 0)
    }
    return result, nil
}
