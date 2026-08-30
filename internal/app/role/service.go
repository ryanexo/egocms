package role

import (
    "context"
    "errors"
    "slices"
    "sync"
    
    "cms/internal/app/role/api"
    "cms/internal/app/role/errno"
    "cms/internal/app/role/model"
    "cms/internal/infra/casbin"
    "cms/internal/infra/store/modeltype"
    "cms/internal/public/apitype"
    
    "gorm.io/gen"
)

type Repo interface {
    Create(ctx context.Context, role *model.Role) error
    Update(ctx context.Context, role *model.Role) (gen.ResultInfo, error)
    Delete(ctx context.Context, id uint64) error
    FindByID(ctx context.Context, id uint64) (*model.Role, error)
    List(ctx context.Context, params *api.RoleListParams) ([]*model.Role, int64, error)
}

type Service struct {
    repo   Repo
    casbin *casbin.RoleCasbin
    mu     sync.Mutex
}

func NewRoleService(repo Repo, roleCasbin *casbin.RoleCasbin) *Service {
    return &Service{repo: repo, casbin: roleCasbin}
}

func (s *Service) Create(ctx context.Context, params api.RoleCreateParams) (uint64, error) {
    role, err := buildRole(params.Name, params.Description)
    if err != nil {
        return 0, err
    }
    inheritList := uniqueRoleParams(params.InheritList)
    if len(inheritList) > model.RoleInheritanceLimit {
        return 0, errno.ErrTooManyInheritedRoles
    }
    
    s.mu.Lock()
    defer s.mu.Unlock()
    
    if err = s.canInherit(ctx, 0, inheritList); err != nil {
        return 0, err
    }
    if err = s.repo.Create(ctx, role); err != nil {
        return 0, err
    }
    if err = s.reinherit(role.ID, casbinRoleNames(inheritList)); err != nil {
        rollbackPolicyErr := s.reinherit(role.ID, nil)
        rollbackRoleErr := s.repo.Delete(ctx, role.ID)
        return 0, errors.Join(err, rollbackPolicyErr, rollbackRoleErr)
    }
    return role.ID, nil
}

func (s *Service) Update(ctx context.Context, params api.RoleUpdateParams) error {
    updated, err := buildRole(params.Name, params.Description)
    if err != nil {
        return err
    }
    updated.Base = modeltype.Base{ID: params.ID.Uint64()}
    inheritList := uniqueRoleParams(params.InheritList)
    if len(inheritList) > model.RoleInheritanceLimit {
        return errno.ErrTooManyInheritedRoles
    }
    
    s.mu.Lock()
    defer s.mu.Unlock()
    
    current, err := s.repo.FindByID(ctx, params.ID.Uint64())
    if err != nil {
        return err
    }
    if err = s.canInherit(ctx, updated.ID, inheritList); err != nil {
        return err
    }
    
    previousInheritance, err := s.casbin.GetRolesForUser(params.ID.String())
    if err != nil {
        return err
    }
    
    if _, err = s.repo.Update(ctx, updated); err != nil {
        return err
    }
    
    if err = s.reinherit(updated.ID, casbinRoleNames(inheritList)); err != nil {
        _, rollbackRoleErr := s.repo.Update(ctx, current)
        rollbackPolicyErr := s.reinherit(updated.ID, previousInheritance)
        return errors.Join(err, rollbackRoleErr, rollbackPolicyErr)
    }
    return nil
}

func (s *Service) Delete(ctx context.Context, id uint64) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    current, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return err
    }
    if err := s.repo.Delete(ctx, id); err != nil {
        return err
    }
    if _, err = s.casbin.DeleteRole(casbinRoleName(id)); err != nil {
        return err
    }
    if err = s.repo.Create(ctx, current); err != nil {
        return err
    }
    return nil
}

func (s *Service) FindByID(ctx context.Context, id uint64) (*model.Role, error) {
    return s.repo.FindByID(ctx, id)
}

func (s *Service) List(ctx context.Context, params api.RoleListParams) (*apitype.PaginatedResult[*model.Role], error) {
    roles, total, err := s.repo.List(ctx, &params)
    if err != nil {
        return nil, err
    }
    return &apitype.PaginatedResult[*model.Role]{
        Pagination: params.Pagination,
        Total:      total,
        List:       roles,
    }, nil
}

func (s *Service) canInherit(ctx context.Context, roleID uint64, inheritList []uint64) error {
    currentRole := casbinRoleName(roleID)
    
    if slices.Contains(inheritList, roleID) {
        return errno.ErrInheritSelf
    }
    
    for _, inheritedID := range inheritList {
        if roleID == 0 {
            continue
        }
        
        isCircular, err := s.casbin.HasRoleForUser(casbinRoleName(inheritedID), currentRole)
        if err != nil {
            return err
        }
        
        if isCircular {
            if role, err := s.repo.FindByID(ctx, roleID); err != nil {
                return err
            } else {
                return errno.ErrCircular.Format(role.Name)
            }
        }
    }
    return nil
}

func (s *Service) reinherit(roleID uint64, roles []string) error {
    role := casbinRoleName(roleID)
    if _, err := s.casbin.DeleteRolesForUser(role); err != nil {
        return err
    }
    if len(roles) == 0 {
        return nil
    }
    _, err := s.casbin.AddRolesForUser(role, roles)
    return err
}
