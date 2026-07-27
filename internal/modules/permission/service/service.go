package service

import (
    "context"
    
    `cms/internal/infra/persistence/gorm/model`
    "cms/internal/pkg/datatype"
    
    menu "cms/internal/modules/menu/contract"
    "cms/internal/modules/permission/contract"
    "cms/internal/modules/permission/internal/assembler"
    "cms/internal/modules/permission/internal/dto"
)

type PermissionService struct {
    permRepo contract.PermissionRepo
    menuRepo menu.MenuRepo
}

func NewPermissionService(repo contract.PermissionRepo) *PermissionService {
    return &PermissionService{permRepo: repo}
}

func (s *PermissionService) Create(ctx context.Context, params dto.PermissionCreateParams) (datatype.SafeUint64, error) {
    data := assembler.ToPermissionCreateCommand(&params)
    err := s.permRepo.Create(ctx, data)
    if err != nil {
        return 0, err
    }
    return data.ID, nil
}

func (s *PermissionService) Update(ctx context.Context, params dto.PermissionUpdateParams) error {
    data := assembler.ToPermissionUpdateCommand(&params)
    _, err := s.permRepo.Update(ctx, data)
    return err
}

func (s *PermissionService) Delete(ctx context.Context, id datatype.SafeUint64) error {
    _, err := s.permRepo.Delete(ctx, id)
    return err
}

func (s *PermissionService) FindByID(ctx context.Context, id datatype.SafeUint64) (*model.Permission, error) {
    return s.permRepo.FindByID(ctx, id)
}

func (s *PermissionService) FindByMenuID(ctx context.Context, menuID datatype.SafeUint64) ([]*model.Permission, error) {
    return s.permRepo.FindByMenuID(ctx, menuID)
}

func (s *PermissionService) FindAll(ctx context.Context) ([]*model.Permission, error) {
    return s.permRepo.FindAll(ctx)
}

func (s *PermissionService) FindNoMenuID(ctx context.Context) ([]*model.Permission, error) {
    return s.permRepo.FindNoMenuID(ctx)
}

func (s *PermissionService) FindByResource(ctx context.Context, resource string) ([]*model.Permission, error) {
    return s.permRepo.FindByResource(ctx, resource)
}
