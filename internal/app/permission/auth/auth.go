package auth

import (
    `context`
    
    `dpcms/internal/infra/rbac`
    `dpcms/internal/middleware/authz`
)

type permissionChecker struct {
    casbin *rbac.RoleCasbin
}

var _ authz.PermissionChecker = (*permissionChecker)(nil)

func (s permissionChecker) Check(_ context.Context, subject string, object string, action string) (bool, error) {
    return s.casbin.Enforce(subject, object, action)
}

func NewPermissionChecker(casbin *rbac.RoleCasbin) authz.PermissionChecker {
    return &permissionChecker{casbin: casbin}
}
