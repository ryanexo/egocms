package service

import (
    `context`
    
    configkeys `dpcms/internal/app/setting/constant`
    `dpcms/internal/config`
    `dpcms/internal/infra/cache`
    `dpcms/internal/infra/logger`
    `dpcms/internal/infra/persistence/contract`
    `dpcms/internal/infra/persistence/model`
)

type SettingService struct {
    txManager contract.TxManager
    repo      SettingRepo
    cache     *cache.ConfigCache
    logger    *logger.Logger
    cfg       *config.Config
}

func NewSettingService(
    txManager contract.TxManager,
    repo SettingRepo, cache *cache.ConfigCache,
    logger *logger.Logger,
    cfg *config.Config,
) *SettingService {
    srv := &SettingService{
        txManager: txManager,
        repo:      repo,
        cache:     cache,
        logger:    logger,
        cfg:       cfg,
    }
    srv.preload()
    
    return srv
}

func (s SettingService) Add(ctx context.Context, key string, value string) error {
    err := s.repo.Add(ctx, &model.Config{Key: key, Value: value})
    if err != nil {
        return err
    }
    s.cache.Set(key, value, 0)
    return nil
}

func (s SettingService) Update(ctx context.Context, key string, value string) error {
    _, err := s.repo.Update(ctx, &model.Config{Key: key, Value: value})
    if err != nil {
        return err
    }
    s.cache.Set(key, value, 0)
    return err
}

func (s SettingService) Remove(ctx context.Context, key string) error {
    _, err := s.repo.Remove(ctx, key)
    if err != nil {
        return err
    }
    s.cache.Del(key)
    return err
}

func (s SettingService) Get(key string) (string, bool) {
    return s.cache.Get(key)
}

func (s SettingService) preload() {
    data, err := s.repo.GetAll(context.Background())
    if err != nil {
        s.logger.App.Error(err.Error())
        return
    }
    for _, v := range data {
        s.cache.Set(v.Key, v.Value, 0)
    }
    s.cache.Set(configkeys.FileDriver, s.cfg.File.Default, 0)
}
