package service

import (
    `context`
    
    `dpcms/internal/app/permission/internal/assembler`
    `dpcms/internal/app/permission/internal/dto`
    `dpcms/internal/infra`
    `dpcms/internal/infra/persistence/datatype`
)

type PermissionService struct {
    Infra *infra.Infra
}

func NewService(infra *infra.Infra) (*PermissionService, error) {
    return &PermissionService{infra}, nil
}

func (s *PermissionService) Create(ctx context.Context, params dto.PermissionCreateParams) (datatype.SafeUint64, error) {
    data := assembler.BuildPermissionCreateCommand(&params)
    
    return data.ID, nil
}
