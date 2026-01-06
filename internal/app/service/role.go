package service

import (
    `context`
    
    roleAssembler `dpcms/internal/app/assembler/role`
    `dpcms/internal/app/dto`
    `dpcms/internal/app/erroz`
    `dpcms/internal/app/util/rbacutil`
    `dpcms/internal/infra`
    `dpcms/internal/infra/persistence/datatype`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/infra/persistence/query`
    `dpcms/internal/infra/persistence/repo`
)

type Role struct {
    rbac    *RBAC
    persist *query.Query
}

func NewRoleService(i *infra.Infra, rbac *RBAC) (*Role, error) {
    return &Role{rbac: rbac, persist: i.Query}, nil
}

func (srv Role) Create(ctx context.Context, params dto.RoleCreateParams) (datatype.SafeUint64, error) {
    data := &model.Role{Name: params.Name, Description: params.Description}
    err := srv.persist.Transaction(func(tx *query.Query) error {
        roleRepo := repo.NewRoleRepo(srv.persist)
        txErr := roleRepo.Create(ctx, data)
        if txErr != nil {
            return txErr
        }
        
        if len(params.InheritList) == 0 {
            return nil
        }
        inheritList := make([]string, 0, len(params.InheritList))
        
        for _, inheritID := range params.InheritList {
            inheritList = append(inheritList, rbacutil.GetRoleSubject(inheritID.Raw()))
        }
        
        enforcer := srv.rbac.GetEnforcer()
        _, txErr = enforcer.AddRolesForUser(rbacutil.GetRoleSubject(data.ID.Raw()), inheritList)
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

func (srv Role) Update(ctx context.Context, params dto.RoleUpdateParams) error {
    return srv.persist.Transaction(func(tx *query.Query) error {
        roleRepo := repo.NewRoleRepo(srv.persist)
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
        
        enforcer := srv.rbac.GetEnforcer()
        currentRoleID := rbacutil.GetRoleSubject(params.ID.Raw())
        
        for _, roleID := range params.InheritList {
            inheritRoleID := rbacutil.GetRoleSubject(roleID.Raw())
            isCircular, err := enforcer.HasRoleForUser(inheritRoleID, currentRoleID)
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
            
            _, err = enforcer.AddRoleForUser(currentRoleID, rbacutil.GetRoleSubject(roleID.Raw()))
            if err != nil {
                return err
            }
        }
        
        return enforcer.SavePolicy()
    })
}

func (srv Role) Delete(ctx context.Context, id datatype.SafeUint64) error {
    return srv.persist.Transaction(func(tx *query.Query) error {
        _, err := repo.NewRoleRepo(tx).Delete(ctx, id)
        if err != nil {
            return err
        }
        
        enforcer := srv.rbac.GetEnforcer()
        sub := rbacutil.GetRoleSubject(id.Raw())
        _, err = enforcer.DeleteRole(sub)
        if err != nil {
            return err
        }
        
        return enforcer.SavePolicy()
    })
}

func (srv Role) FindByID(ctx context.Context, id datatype.SafeUint64) (*dto.Role, error) {
    data, err := repo.NewRoleRepo(srv.persist).FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    return roleAssembler.BuildRoleDTO(data), nil
}

func (srv Role) List(ctx context.Context, params dto.RoleListParams) (*dto.PaginatedResult[*dto.Role], error) {
    data, total, err := repo.NewRoleRepo(srv.persist).List(ctx, &params)
    if err != nil {
        return nil, err
    }
    return &dto.PaginatedResult[*dto.Role]{
        Total:    total,
        PageNo:   params.PageNo,
        PageSize: params.PageSize,
        List:     roleAssembler.BuildRoleListDTO(data),
    }, nil
}
