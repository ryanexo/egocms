package adapter

import (
    `context`
    
    `cms/internal/infra/casbin`
    `cms/internal/middleware/authz`
)

type permissionChecker struct {
    casbin *casbin.RoleCasbin
}

var _ authz.PermissionChecker = (*permissionChecker)(nil)

func (s permissionChecker) Check(_ context.Context, subject string, object string, action string) (bool, error) {
    return s.casbin.Enforce(subject, object, action)
}

func NewPermissionChecker(casbin *casbin.RoleCasbin) authz.PermissionChecker {
    return &permissionChecker{casbin: casbin}
}
