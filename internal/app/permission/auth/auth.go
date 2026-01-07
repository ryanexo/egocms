package auth

import (
    `context`
    
    `dpcms/internal/middleware/authz`
    
    `github.com/casbin/casbin/v2`
)

type permissionChecker struct {
    casbin *casbin.Enforcer
}

var _ authz.PermissionChecker = (*permissionChecker)(nil)

func (s permissionChecker) Check(_ context.Context, subject string, object string, action string) (bool, error) {
    return s.casbin.Enforce(subject, object, action)
}

func NewPermissionChecker(casbin *casbin.Enforcer) authz.PermissionChecker {
    return &permissionChecker{casbin: casbin}
}
