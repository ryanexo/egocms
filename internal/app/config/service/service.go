package service

import (
    `context`
    
    `dpcms/internal/app/config/internal/dto`
    `dpcms/internal/infra/persistence/contract`
    `dpcms/internal/infra/persistence/model`
)

type ConfigService struct {
    TxManager contract.TxManager
    Repo      ConfigRepo
}

func (s ConfigService) Add(ctx context.Context, key string, value string) error {
    return s.Repo.Add(ctx, &model.Config{Key: key, Value: value})
}

func (s ConfigService) Update(ctx context.Context, key string, value string) error {
    _, err := s.Repo.Update(ctx, &model.Config{Key: key, Value: value})
    return err
}

func (s ConfigService) Remove(ctx context.Context, key string) error {
    _, err := s.Repo.Remove(ctx, key)
    return err
}

func (s ConfigService) Get(ctx context.Context, key string) (string, error) {
    return s.Repo.Get(ctx, key)
}

func (s ConfigService) GetAll(ctx context.Context) (dto.ConfigMap, error) {
    data, err := s.Repo.GetAll(ctx)
    if err != nil {
        return nil, err
    }
    result := make(dto.ConfigMap, len(data))
    for _, v := range data {
        result[v.Key] = v.Value
    }
    return result, nil
}
