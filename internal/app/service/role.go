package service

import (
    `context`
    
    `dpcms/internal/app/dto`
    `dpcms/internal/app/dto/type`
    `dpcms/internal/app/erroz`
    `dpcms/internal/app/util/rbacutil`
    `dpcms/internal/infra`
    `dpcms/internal/infra/persistence/datatype`
    `dpcms/internal/infra/persistence/dbscope`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/infra/persistence/query`
    
    `github.com/jinzhu/copier`
)

type Role struct {
    rbac    *RBAC
    persist *query.Query
}

func NewRoleService(i *infra.Infra, rbac *RBAC) (*Role, error) {
    return &Role{rbac: rbac, persist: i.Query}, nil
}

func (srv Role) Create(ctx context.Context, params dto.RoleCreateParams) (result *dto.Role, err error) {
    roleModel := model.Role{Name: params.Name, Description: params.Description}
    
    err = srv.persist.Transaction(func(tx *query.Query) error {
        queryCtx := tx.WithContext(ctx)
        txErr := queryCtx.Role.Create(&roleModel)
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
        _, txErr = enforcer.AddRolesForUser(rbacutil.GetRoleSubject(roleModel.ID.Raw()), inheritList)
        if txErr != nil {
            return txErr
        }
        
        txErr = enforcer.SavePolicy()
        if txErr != nil {
            return txErr
        }
        
        return nil
    })
    if err != nil {
        return
    }
    
    inheritList := make([]uint64, 0, len(params.InheritList))
    for _, inheritId := range params.InheritList {
        inheritList = append(inheritList, inheritId.Raw())
    }
    
    roles, err := srv.persist.WithContext(ctx).Role.Where(srv.persist.Role.ID.In(inheritList...)).Find()
    if err != nil {
        return
    }
    err = copier.Copy(&result, &roleModel)
    if err != nil {
        return
    }
    err = copier.Copy(&result.InheritList, &roles)
    if err != nil {
        return
    }
    
    return
}

func (srv Role) Update(ctx context.Context, params dto.RoleUpdateParams) error {
    return srv.persist.Transaction(func(tx *query.Query) error {
        queryCtx := tx.WithContext(ctx)
        
        _, err := queryCtx.Role.Where(srv.persist.Role.ID.Eq(params.ID.Raw())).Updates(model.Role{
            Name:        params.Name,
            Description: params.Description,
        })
        if err != nil {
            return err
        }
        
        if len(params.InheritList) == 0 {
            return nil
        }
        
        enforcer := srv.rbac.GetEnforcer()
        currentRoleName := rbacutil.GetRoleSubject(params.ID.Raw())
        for _, inheritRoleID := range params.InheritList {
            inheritRoleName := rbacutil.GetRoleSubject(inheritRoleID.Raw())
            isCircleRelate, err := enforcer.HasRoleForUser(inheritRoleName, currentRoleName)
            if err != nil {
                return err
            }
            
            if isCircleRelate {
                childRole, err := tx.Role.WithContext(ctx).Where(tx.Role.ID.Eq(inheritRoleID.Raw())).First()
                if err != nil {
                    return err
                }
                return erroz.RoleCircular.Format(childRole.Name).ToError()
            }
            
            _, err = enforcer.AddRoleForUser(currentRoleName, rbacutil.GetRoleSubject(inheritRoleID.Raw()))
            if err != nil {
                return err
            }
        }
        
        err = enforcer.SavePolicy()
        if err != nil {
            return err
        }
        
        return nil
    })
}

func (srv Role) Delete(ctx context.Context, id datatype.SafeUint64) error {
    return srv.persist.Transaction(func(tx *query.Query) error {
        queryCtx := tx.WithContext(ctx)
        _, err := queryCtx.Role.Where(tx.Role.ID.Eq(id.Raw())).Delete()
        if err != nil {
            return err
        }
        enforcer := srv.rbac.GetEnforcer()
        _, err = enforcer.DeleteRole(rbacutil.GetRoleSubject(id.Raw()))
        if err != nil {
            return err
        }
        err = enforcer.SavePolicy()
        if err != nil {
            return err
        }
        return nil
    })
}

func (srv Role) FindRoleByID(ctx context.Context, id datatype.SafeUint64) (result dto.Role, err error) {
    dao := srv.persist.Role
    r, err := dao.WithContext(ctx).Where(dao.ID.Eq(id.Raw())).First()
    if err != nil {
        return
    }
    
    err = copier.Copy(&result, &r)
    if err != nil {
        return
    }
    
    currentRoleName := rbacutil.GetRoleSubject(r.ID.Raw())
    inheritRoleNames, err := srv.rbac.GetEnforcer().GetRolesForUser(currentRoleName)
    if err != nil {
        return
    }
    
    idList, err := rbacutil.ParseRoleSubject(inheritRoleNames...)
    inheritRoles, err := dao.WithContext(ctx).Where(dao.ID.In(idList...)).Find()
    if err != nil {
        return
    }
    
    err = copier.Copy(&result.InheritList, &inheritRoles)
    return
}

func (srv Role) List(ctx context.Context, params dto.RoleListParams) (*dtotype.PaginatedResult[*model.Role], error) {
    dao := srv.persist.Role
    q := dao.WithContext(ctx).Scopes(dbscope.Paginate(params.PageNo, params.PageSize))
    
    if params.Name != nil {
        q = q.Where(dao.Name.Eq(*params.Name))
    }
    if params.Description != nil {
        q = q.Where(dao.Description.Eq(*params.Description))
    }
    if params.InheritId != nil {
        sub := rbacutil.GetRoleSubject(params.InheritId.Raw())
        inheritSubjects, err := srv.rbac.GetEnforcer().GetUsersForRole(sub)
        if err != nil {
            return nil, err
        }
        inheritIdList, err := rbacutil.ParseRoleSubject(inheritSubjects...)
        if err != nil {
            return nil, err
        }
        q = q.Where(dao.ID.In(inheritIdList...))
    }
    
    count, err := q.Count()
    if err != nil {
        return nil, err
    }
    roles, err := q.Find()
    if err != nil {
        return nil, err
    }
    
    return &dtotype.PaginatedResult[*model.Role]{
        Total:    count,
        PageNo:   params.PageNo,
        PageSize: params.PageSize,
        List:     roles,
    }, nil
}
