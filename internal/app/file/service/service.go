package service

import (
    `context`
    
    `dpcms/internal/infra/persistence/query`
)

type FileService struct {
    Query *query.Query
}

func (s FileService) Save(ctx context.Context, path string) {}
