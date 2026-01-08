package service

import (
    `context`
    
    `dpcms/internal/app/permission/internal/assembler`
    `dpcms/internal/app/permission/internal/dto`
    `dpcms/internal/infra/persistence/datatype`
    `dpcms/internal/infra/persistence/query`
)

type PermissionService struct {
    Query *query.Query
}

func (s *PermissionService) Create(ctx context.Context, params dto.PermissionCreateParams) (datatype.SafeUint64, error) {
    data := assembler.BuildPermissionCreateCommand(&params)
    
    return data.ID, nil
}
