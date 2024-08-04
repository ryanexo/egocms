package service

import (
    "context"
    
    `dpcms/api/repository`
    `dpcms/model`
)

type Category struct {
    repo repository.Repositories
}

func (srv Category) Create(ctx context.Context, category *model.Category) error {
    
    return nil
}

func (srv Category) Update(ctx context.Context, category *model.Category) error {
    return nil
}
