package service

import (
    `context`
    
    `cms/internal/config`
    contract2 `cms/internal/domain/setting/contract`
    `cms/internal/infra/logger`
    `cms/internal/infra/persistence/contract`
    `cms/internal/infra/persistence/model`
)

type SettingService struct {
    txManager contract.Transactor
    repo      contract2.SettingRepo
    logger    *logger.Logger
    cfg       *config.Config
}

func NewSettingService(
    txManager contract.Transactor,
    repo contract2.SettingRepo,
    logger *logger.Logger,
    cfg *config.Config,
) *SettingService {
    return &SettingService{
        txManager: txManager,
        repo:      repo,
        logger:    logger,
        cfg:       cfg,
    }
}

func (s SettingService) Add(ctx context.Context, key string, value string) error {
    return s.repo.Add(ctx, &model.Setting{Key: key, Value: value})
}

func (s SettingService) Update(ctx context.Context, key string, value string) error {
    _, err := s.repo.Update(ctx, &model.Setting{Key: key, Value: value})
    return err
}

func (s SettingService) Remove(ctx context.Context, key string) error {
    _, err := s.repo.Remove(ctx, key)
    return err
}

func (s SettingService) Get(ctx context.Context, key string) (string, error) {
    return s.repo.Get(ctx, key)
}
