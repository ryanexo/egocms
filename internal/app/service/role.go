package service

import (
    `context`
    
    `dpcms/internal/app/dto`
    `dpcms/internal/app/erroz`
    `dpcms/internal/app/helper/rbachelper`
    `dpcms/internal/app/service/internal/common`
    `dpcms/internal/infra`
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
            inheritList = append(inheritList, rbachelper.GetRoleSubject(inheritID))
        }
        
        enforcer := srv.rbac.GetEnforcer()
        _, txErr = enforcer.AddRolesForUser(rbachelper.GetRoleSubject(roleModel.ID), inheritList)
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
    
    roles, err := srv.persist.WithContext(ctx).Role.Where(srv.persist.Role.ID.In(params.InheritList...)).Find()
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
        
        _, err := queryCtx.Role.Where(srv.persist.Role.ID.Eq(params.ID)).Updates(model.Role{
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
        currentRoleName := rbachelper.GetRoleSubject(params.ID)
        for _, inheritRoleID := range params.InheritList {
            inheritRoleName := rbachelper.GetRoleSubject(inheritRoleID)
            isCircleRelate, err := enforcer.HasRoleForUser(inheritRoleName, currentRoleName)
            if err != nil {
                return err
            }
            
            if isCircleRelate {
                childRole, err := srv.persist.Role.WithContext(ctx).Where(srv.persist.Role.ID.Eq(inheritRoleID)).First()
                if err != nil {
                    return err
                }
                return erroz.RoleCircularReference.Format(childRole.Name).ToError()
            }
            
            _, err = enforcer.AddRoleForUser(currentRoleName, rbachelper.GetRoleSubject(inheritRoleID))
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

func (srv Role) Delete(ctx context.Context, id uint64) error {
    return srv.persist.Transaction(func(tx *query.Query) error {
        queryCtx := tx.WithContext(ctx)
        _, err := queryCtx.Role.Where(srv.persist.Role.ID.Eq(id)).Delete()
        if err != nil {
            return err
        }
        enforcer := srv.rbac.GetEnforcer()
        _, err = enforcer.DeleteRole(rbachelper.GetRoleSubject(id))
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

func (srv Role) FindRoleById(ctx context.Context, id uint64) (result dto.Role, err error) {
    dao := srv.persist.Role
    r, err := dao.WithContext(ctx).Where(dao.ID.Eq(id)).First()
    if err != nil {
        return
    }
    
    err = copier.Copy(&result, &r)
    if err != nil {
        return
    }
    
    currentRoleName := rbachelper.GetRoleSubject(r.ID)
    inheritRoleNames, err := srv.rbac.GetEnforcer().GetRolesForUser(currentRoleName)
    if err != nil {
        return
    }
    
    idList, err := rbachelper.ParseRoleSubject(inheritRoleNames...)
    inheritRoles, err := dao.WithContext(ctx).Where(dao.ID.In(idList...)).Find()
    if err != nil {
        return
    }
    
    err = copier.Copy(&result.InheritList, &inheritRoles)
    return
}

func (srv Role) List(ctx context.Context, params dto.RoleListParams) (*common.PaginatedResult[*model.Role], error) {
    dao := srv.persist.Role
    q := dao.WithContext(ctx).Scopes(dbscope.Paginate(params.PageNo, params.PageSize))
    
    if params.Name != nil {
        q = q.Where(dao.Name.Eq(*params.Name))
    }
    if params.Description != nil {
        q = q.Where(dao.Description.Eq(*params.Description))
    }
    if params.InheritId != nil {
        sub := rbachelper.GetRoleSubject(*params.InheritId)
        inheritSubjects, err := srv.rbac.GetEnforcer().GetUsersForRole(sub)
        if err != nil {
            return nil, err
        }
        inheritIdList, err := rbachelper.ParseRoleSubject(inheritSubjects...)
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
    
    return &common.PaginatedResult[*model.Role]{
        Total:    count,
        PageNo:   params.PageNo,
        PageSize: params.PageSize,
        List:     roles,
    }, nil
}
