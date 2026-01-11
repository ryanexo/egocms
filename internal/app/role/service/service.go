package service

import (
    `context`
    
    `dpcms/internal/app/role/errno`
    roleAssembler `dpcms/internal/app/role/internal/assembler`
    `dpcms/internal/app/role/internal/dto`
    `dpcms/internal/infra/persistence/contract`
    `dpcms/internal/infra/persistence/datatype`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/infra/persistence/query`
    `dpcms/internal/infra/rbac`
    `dpcms/internal/types`
)

type RoleService struct {
    txManager contract.TxManager
    casbin    *rbac.RoleCasbin
    repo      RoleRepo
}

func NewRoleService(txManager contract.TxManager, casbin *rbac.RoleCasbin, repo RoleRepo) *RoleService {
    return &RoleService{
        txManager: txManager,
        casbin:    casbin,
        repo:      repo,
    }
}

func (s RoleService) Create(ctx context.Context, params dto.RoleCreateParams) (datatype.SafeUint64, error) {
    data := &model.Role{Name: params.Name, Description: params.Description}
    err := s.txManager.Transaction(func(tx *query.Query) error {
        roleRepo := s.repo.CloneWithQuery(tx)
        txErr := roleRepo.Create(ctx, data)
        if txErr != nil {
            return txErr
        }
        
        if len(params.InheritList) == 0 {
            return nil
        }
        inheritList := make([]string, 0, len(params.InheritList))
        
        for _, inheritID := range params.InheritList {
            inheritList = append(inheritList, inheritID.String())
        }
        
        enforcer := s.casbin
        _, txErr = enforcer.AddRolesForUser(data.ID.String(), inheritList)
        if txErr != nil {
            return txErr
        }
        
        txErr = enforcer.SavePolicy()
        if txErr != nil {
            return txErr
        }
        
        return nil
    })
    return data.ID, err
}

func (s RoleService) Update(ctx context.Context, params dto.RoleUpdateParams) error {
    return s.txManager.Transaction(func(tx *query.Query) error {
        roleRepo := s.repo.CloneWithQuery(tx)
        data := &model.Role{
            Name:        params.Name,
            Description: params.Description,
        }
        
        _, err := roleRepo.Update(ctx, data)
        if err != nil {
            return err
        }
        
        if len(params.InheritList) == 0 {
            return nil
        }
        
        enforcer := s.casbin
        currentRoleID := params.ID.String()
        
        for _, roleID := range params.InheritList {
            isCircular, err := enforcer.HasRoleForUser(roleID.String(), currentRoleID)
            if err != nil {
                return err
            }
            
            if isCircular {
                role, err := roleRepo.FindByID(ctx, roleID)
                if err != nil {
                    return err
                }
                return errno.RoleCircular.Format(role.Name).ToError()
            }
            
            _, err = enforcer.AddRoleForUser(currentRoleID, roleID.String())
            if err != nil {
                return err
            }
        }
        
        return enforcer.SavePolicy()
    })
}

func (s RoleService) Delete(ctx context.Context, id datatype.SafeUint64) error {
    return s.txManager.Transaction(func(tx *query.Query) error {
        _, err := s.repo.CloneWithQuery(tx).Delete(ctx, id)
        if err != nil {
            return err
        }
        
        enforcer := s.casbin
        _, err = enforcer.DeleteRole(id.String())
        if err != nil {
            return err
        }
        
        return enforcer.SavePolicy()
    })
}

func (s RoleService) FindByID(ctx context.Context, id datatype.SafeUint64) (*dto.Role, error) {
    data, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    return roleAssembler.BuildRoleDTO(data), nil
}

func (s RoleService) List(ctx context.Context, params dto.RoleListParams) (*types.PaginatedResult[*dto.Role], error) {
    data, total, err := s.repo.List(ctx, &params)
    if err != nil {
        return nil, err
    }
    return &types.PaginatedResult[*dto.Role]{
        Pagination: params.Pagination,
        Total:      total,
        List:       roleAssembler.BuildRoleListDTO(data),
    }, nil
}
