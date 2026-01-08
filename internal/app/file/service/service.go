package service

import `dpcms/internal/infra/persistence/query`

type FileService struct {
    Query *query.Query
}

func (s FileService) Upload() {}
