package service

import (
    `context`
    
    roleAssembler `dpcms/internal/app/role/internal/assembler`
    `dpcms/internal/app/role/internal/dto`
    `dpcms/internal/app/role/repo`
    `dpcms/internal/erroz`
    `dpcms/internal/infra`
    `dpcms/internal/infra/persistence/datatype`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/infra/persistence/query`
    `dpcms/internal/types`
)

type RoleService struct {
    persist *query.Query
    infra   *infra.Infra
}

func NewRoleService(persist *query.Query, infra *infra.Infra) (*RoleService, error) {
    return &RoleService{persist, infra}, nil
}

func (s RoleService) Create(ctx context.Context, params dto.RoleCreateParams) (datatype.SafeUint64, error) {
    data := &model.Role{Name: params.Name, Description: params.Description}
    err := s.persist.Transaction(func(tx *query.Query) error {
        roleRepo := repo.NewRoleRepo(s.persist)
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
        
        enforcer := s.infra.Casbin
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
    return s.persist.Transaction(func(tx *query.Query) error {
        roleRepo := repo.NewRoleRepo(s.persist)
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
        
        enforcer := s.infra.Casbin
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
                return erroz.RoleCircular.Format(role.Name).ToError()
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
    return s.persist.Transaction(func(tx *query.Query) error {
        _, err := repo.NewRoleRepo(tx).Delete(ctx, id)
        if err != nil {
            return err
        }
        
        enforcer := s.infra.Casbin
        _, err = enforcer.DeleteRole(id.String())
        if err != nil {
            return err
        }
        
        return enforcer.SavePolicy()
    })
}

func (s RoleService) FindByID(ctx context.Context, id datatype.SafeUint64) (*dto.Role, error) {
    data, err := repo.NewRoleRepo(s.persist).FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    return roleAssembler.BuildRoleDTO(data), nil
}

func (s RoleService) List(ctx context.Context, params dto.RoleListParams) (*types.PaginatedResult[*dto.Role], error) {
    data, total, err := repo.NewRoleRepo(s.persist).List(ctx, &params)
    if err != nil {
        return nil, err
    }
    return &types.PaginatedResult[*dto.Role]{
        Pagination: params.Pagination,
        Total:      total,
        List:       roleAssembler.BuildRoleListDTO(data),
    }, nil
}
