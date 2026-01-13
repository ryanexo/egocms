package service

import (
    `context`
    
    `dpcms/internal/app/permission/internal/assembler`
    `dpcms/internal/app/permission/internal/dto`
    `dpcms/internal/infra/persistence/datatype`
    `dpcms/internal/infra/persistence/model`
)

type PermissionService struct {
    repo PermissionRepo
}

func NewPermissionService(repo PermissionRepo) *PermissionService {
    return &PermissionService{repo: repo}
}

func (s *PermissionService) Create(ctx context.Context, params dto.PermissionCreateParams) (datatype.SafeUint64, error) {
    data := assembler.ToPermissionCreateCommand(&params)
    err := s.repo.Create(ctx, data)
    if err != nil {
        return 0, err
    }
    return data.ID, nil
}

func (s *PermissionService) Update(ctx context.Context, params dto.PermissionUpdateParams) error {
    data := assembler.ToPermissionUpdateCommand(&params)
    _, err := s.repo.Update(ctx, data)
    return err
}

func (s *PermissionService) Delete(ctx context.Context, id datatype.SafeUint64) error {
    _, err := s.repo.Delete(ctx, id)
    return err
}

func (s *PermissionService) FindByID(ctx context.Context, id datatype.SafeUint64) (*model.Permission, error) {
    return s.repo.FindByID(ctx, id)
}

func (s *PermissionService) FindByMenuID(ctx context.Context, menuID datatype.SafeUint64) ([]*model.Permission, error) {
    return s.repo.FindByMenuID(ctx, menuID)
}

func (s *PermissionService) FindAll(ctx context.Context) ([]*model.Permission, error) {
    return s.repo.FindAll(ctx)
}

func (s *PermissionService) FindNoMenuID(ctx context.Context) ([]*model.Permission, error) {
    return s.repo.FindNoMenuID(ctx)
}

func (s *PermissionService) FindByResource(ctx context.Context, resource string) ([]*model.Permission, error) {
    return s.repo.FindByResource(ctx, resource)
}
