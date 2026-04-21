package domain

import (
    `context`
    
    `cms/internal/infra/persist/model`
)

type CategoryDomainService struct {
    repo CategoryRepo
}

func NewCategoryDomainService(repo CategoryRepo) *CategoryDomainService {
    return &CategoryDomainService{repo: repo}
}

func (s *CategoryDomainService) Create(ctx context.Context, category Category) (*model.Article, error) {
}
