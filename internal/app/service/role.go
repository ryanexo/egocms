package service

import (
    `context`
    
    `dpcms/internal/app/erroz`
    `dpcms/internal/app/helper/dbscope`
    `dpcms/internal/app/helper/rbac`
    `dpcms/internal/app/service/types`
    `dpcms/internal/app/service/types/role`
    `dpcms/internal/database/model`
    `dpcms/internal/database/query`
    `dpcms/internal/infra`
    `dpcms/internal/packages/database`
    
    `github.com/jinzhu/copier`
)

type RoleService struct {
    rbac  RBACService
    query *query.Query
}

func NewRoleService(infra *infra.Infra, rbac RBACService) (RoleService, error) {
    return RoleService{rbac: rbac, query: query.Use(infra.DB)}, nil
}

func (srv RoleService) Create(ctx context.Context, params role.CreateParams) (result role.Detail, err error) {
    roleModel := model.Role{Name: params.Name, Description: params.Description}
    
    err = srv.query.Transaction(func(tx *query.Query) error {
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
            inheritList = append(inheritList, rbac.GetRoleSubject(inheritID))
        }
        
        _, txErr = srv.rbac.AddRolesForUser(rbac.GetRoleSubject(roleModel.ID), inheritList)
        if txErr != nil {
            return txErr
        }
        
        txErr = srv.rbac.SavePolicy()
        if txErr != nil {
            return txErr
        }
        
        return nil
    })
    if err != nil {
        return
    }
    
    roles, err := srv.query.WithContext(ctx).Role.Where(srv.query.Role.ID.In(params.InheritList...)).Find()
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

func (srv RoleService) Update(ctx context.Context, params role.UpdateParams) (result model.Role, err error) {
    err = srv.query.Transaction(func(tx *query.Query) error {
        queryCtx := tx.WithContext(ctx)
        
        _, err := queryCtx.Role.Where(srv.query.Role.ID.Eq(params.ID)).Updates(model.Role{
            Model: database.Model{
                ID: params.ID,
            },
            Name:        params.Name,
            Description: params.Description,
        })
        if err != nil {
            return err
        }
        
        if len(params.InheritList) == 0 {
            return nil
        }
        
        currentRoleName := rbac.GetRoleSubject(params.ID)
        for _, inheritRoleID := range params.InheritList {
            inheritRoleName := rbac.GetRoleSubject(inheritRoleID)
            linked, err := srv.rbac.HasRoleForUser(inheritRoleName, currentRoleName)
            if err != nil {
                return err
            }
            if linked {
                linkedRole, err := srv.query.Role.WithContext(ctx).Where(srv.query.Role.ID.Eq(inheritRoleID)).First()
                if err != nil {
                    return err
                }
                return erroz.ErrRoleCircularReference.Format(linkedRole.Name).ToError()
            }
            _, err = srv.rbac.AddRoleForUser(currentRoleName, rbac.GetRoleSubject(inheritRoleID))
            if err != nil {
                return err
            }
        }
        
        err = srv.rbac.SavePolicy()
        if err != nil {
            return err
        }
        
        return nil
    })
    
    r, err := srv.query.WithContext(ctx).Role.Where(srv.query.Role.ID.Eq(params.ID)).First()
    if err != nil {
        return
    }
    result = *r
    return
}

func (srv RoleService) Delete(ctx context.Context, params role.DeleteParams) error {
    return srv.query.Transaction(func(tx *query.Query) error {
        queryCtx := tx.WithContext(ctx)
        _, err := queryCtx.Role.Where(srv.query.Role.ID.Eq(params.ID)).Delete()
        if err != nil {
            return err
        }
        _, err = srv.rbac.DeleteRole(rbac.GetRoleSubject(params.ID))
        if err != nil {
            return err
        }
        err = srv.rbac.SavePolicy()
        if err != nil {
            return err
        }
        return nil
    })
}

func (srv RoleService) FindByIDWithInherit(ctx context.Context, id int64) (result role.Detail, err error) {
    dao := srv.query.Role
    r, err := dao.WithContext(ctx).Where(dao.ID.Eq(id)).First()
    if err != nil {
        return
    }
    
    err = copier.Copy(&result, &r)
    if err != nil {
        return
    }
    
    currentRoleName := rbac.GetRoleSubject(r.ID)
    inheritRoleNames, err := srv.rbac.GetRolesForUser(currentRoleName)
    if err != nil {
        return
    }
    
    idList, err := rbac.ParseRoleSubject(inheritRoleNames...)
    inheritRoles, err := dao.WithContext(ctx).Where(dao.ID.In(idList...)).Find()
    if err != nil {
        return
    }
    
    err = copier.Copy(&result.InheritList, &inheritRoles)
    return
}

func (srv RoleService) List(ctx context.Context, params role.ListRetrieveParams) (result types.PaginatedResult[role.ListDetail], err error) {
    dao := srv.query.Role
    q := dao.WithContext(ctx).Scopes(dbscope.Paginate(params.PageNo, params.PageSize))
    if params.Name != nil {
        q = q.Where(dao.Name.Eq(*params.Name))
    }
    if params.Description != nil {
        q = q.Where(dao.Description.Eq(*params.Description))
    }
    count, err := q.Count()
    if err != nil {
        return
    }
    result.Total = count
    result.PageNo = params.PageNo
    
    roles, err := q.Find()
    if err != nil {
        return
    }
    err = copier.Copy(&result.List, &roles)
    return
}
