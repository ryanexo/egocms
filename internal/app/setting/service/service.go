package service

import (
	"cms/internal/infra/persistence"
	"context"

	contract2 "cms/internal/app/setting/contract"
	"cms/internal/config"
	"cms/internal/infra/logger"
	"cms/internal/infra/persistence/model"
)

type SettingService struct {
	txManager persistence.Transactor
	repo      contract2.SettingRepo
	logger    *logger.Logger
	cfg       *config.Config
}

func NewSettingService(
	txManager persistence.Transactor,
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
