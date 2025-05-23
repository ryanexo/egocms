package service

import (
    `context`
    
    `dpcms/api/errors/role_error`
    `dpcms/api/infra`
    `dpcms/api/service/internal/rbac`
    `dpcms/model`
    `dpcms/model/query`
    `dpcms/utils/db_scopes`
)

type RoleService struct {
    rbac  RBACService
    query *query.Query
}

func NewRoleService(infra *infra.Infra, rbac RBACService) (*RoleService, error) {
    return &RoleService{rbac: rbac, query: query.Use(infra.DB)}, nil
}

func (srv *RoleService) Create(ctx context.Context, role *model.Role, inheritList []int64) error {
    return srv.query.Transaction(func(tx *query.Query) error {
        queryCtx := tx.WithContext(ctx)
        err := queryCtx.Role.Create(role)
        if err != nil {
            return err
        }
        if len(inheritList) == 0 {
            return nil
        }
        
        for _, inheritID := range inheritList {
            _, err := srv.rbac.AddRolesForUser(rbac.GetRoleSubject(role.ID), []string{rbac.GetRoleSubject(inheritID)})
            if err != nil {
                return err
            }
        }
        
        return nil
    })
}

func (srv *RoleService) Update(ctx context.Context, role *model.Role, inheritList []int64) error {
    return srv.query.Transaction(func(tx *query.Query) error {
        queryCtx := tx.WithContext(ctx)
        
        _, err := queryCtx.Role.Updates(role)
        if err != nil {
            return err
        }
        
        if len(inheritList) == 0 {
            return nil
        }
        
        currentRoleName := rbac.GetRoleSubject(role.ID)
        for _, inheritRoleID := range inheritList {
            inheritRoleName := rbac.GetRoleSubject(inheritRoleID)
            linked, err := srv.rbac.HasRoleForUser(currentRoleName, inheritRoleName)
            if err != nil {
                return err
            }
            if linked {
                linkedRole, err := srv.query.Role.WithContext(ctx).Where(srv.query.Role.ID.Eq(role.ID)).First()
                if err != nil {
                    return err
                }
                return role_error.ErrRoleCircularReference.Format(linkedRole.Name, role.Name).ToError()
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
}

func (srv *RoleService) Delete(ctx context.Context, role *model.Role) error {
    return srv.query.Transaction(func(tx *query.Query) error {
        queryCtx := tx.WithContext(ctx)
        _, err := queryCtx.Role.Delete(role)
        if err != nil {
            return err
        }
        _, err = srv.rbac.DeleteRole(rbac.GetRoleSubject(role.ID))
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

func (srv *RoleService) FindByID(ctx context.Context, id int64) (*model.Role, error) {
    dao := srv.query.Role
    role, err := dao.WithContext(ctx).Where(dao.ID.Eq(id)).First()
    if err != nil {
        return nil, err
    }
    return role, nil
}

func (srv *RoleService) FindByIDWithInherit(ctx context.Context, id int64) (*model.Role, error) {
    role, err := srv.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    currentRoleName := rbac.GetRoleSubject(role.ID)
    inheritRoleNames, err := srv.rbac.GetImplicitRolesForUser(currentRoleName)
    if err != nil {
        return nil, err
    }
    idList, err := rbac.ParseRoleSubject(inheritRoleNames...)
    inheritRoles, err := srv.query.WithContext(ctx).Role.Where(srv.query.Role.ID.In(idList...)).Find()
    if err != nil {
        return nil, err
    }
    role.InheritList = inheritRoles
    return role, nil
}

func (srv *RoleService) List(ctx context.Context, pageNo int, pageSize int) ([]*model.Role, error) {
    return srv.query.WithContext(ctx).Role.Scopes(db_scopes.Paginate(pageNo, pageSize)).Find()
}
