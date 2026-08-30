package setting

import (
	"context"
	"sync"

	"cms/internal/app/setting/model"
	"cms/internal/modules/setting/internal/contract"
)

type Service struct {
	repo  contract.Repo
	cache map[string]string
	mu    sync.RWMutex
}

func NewSettingService(repo contract.Repo) (*Service, error) {
	srv := Service{repo: repo, cache: make(map[string]string)}
	err := srv.reload(context.Background())
	if err != nil {
		return nil, err
	}
	return &srv, nil
}

func (s *Service) Add(ctx context.Context, field, value string) error {
	err := s.repo.Add(ctx, &model.Setting{Field: field, Value: value})
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.cache[field] = value
	s.mu.Unlock()
	return nil
}

func (s *Service) Update(ctx context.Context, field, value string) error {
	_, err := s.repo.Update(ctx, &model.Setting{Field: field, Value: value})
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.cache[field] = value
	s.mu.Unlock()
	return nil
}

func (s *Service) Remove(ctx context.Context, field string) error {
	_, err := s.repo.Remove(ctx, field)
	if err != nil {
		return err
	}
	s.mu.Lock()
	delete(s.cache, field)
	s.mu.Unlock()
	return nil
}

func (s *Service) Get(field string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cache[field]
}

func (s *Service) reload(ctx context.Context) error {
	items, err := s.repo.GetAll(ctx)
	if err != nil {
		return err
	}

	cache := make(map[string]string, len(items))
	for _, item := range items {
		cache[item.Field] = item.Value
	}

	s.mu.Lock()
	clear(s.cache)
	s.cache = cache
	s.mu.Unlock()

	return nil
}
