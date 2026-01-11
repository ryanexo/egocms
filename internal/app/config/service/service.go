package service

import (
    `context`
    `time`
    
    configkeys `dpcms/internal/app/config/constant`
    `dpcms/internal/infra/cache`
    `dpcms/internal/infra/logger`
    `dpcms/internal/infra/persistence/contract`
    `dpcms/internal/infra/persistence/model`
)

type ConfigService struct {
    txManager contract.TxManager
    repo      ConfigRepo
    cache     *cache.ConfigCache
    logger    *logger.Logger
}

func NewConfigService(txManager contract.TxManager, repo ConfigRepo, cache *cache.ConfigCache, logger *logger.Logger) *ConfigService {
    srv := &ConfigService{
        txManager: txManager,
        repo:      repo,
        cache:     cache,
        logger:    logger,
    }
    go srv.backgroundRefresh()
    
    return srv
}

func (s ConfigService) Add(ctx context.Context, key string, value string) error {
    err := s.repo.Add(ctx, &model.Config{Key: key, Value: value})
    if err != nil {
        return err
    }
    s.cache.Set(key, value, 0)
    return nil
}

func (s ConfigService) Update(ctx context.Context, key string, value string) error {
    _, err := s.repo.Update(ctx, &model.Config{Key: key, Value: value})
    if err != nil {
        return err
    }
    s.cache.Set(key, value, 0)
    return err
}

func (s ConfigService) Remove(ctx context.Context, key string) error {
    _, err := s.repo.Remove(ctx, key)
    if err != nil {
        return err
    }
    s.cache.Del(key)
    return err
}

func (s ConfigService) Get(key string) (string, bool) {
    return s.cache.Get(key)
}

func (s ConfigService) refresh() {
    data, err := s.repo.GetAll(context.Background())
    if err != nil {
        s.logger.App.Error(err.Error())
        return
    }
    for _, v := range data {
        s.cache.Set(v.Key, v.Value, 0)
    }
    s.cache.Set(configkeys.FileDriver, "local", 0)
}

func (s ConfigService) backgroundRefresh() {
    s.refresh()
    t := time.NewTicker(time.Minute * 10)
    
    for range t.C {
        s.refresh()
    }
}
