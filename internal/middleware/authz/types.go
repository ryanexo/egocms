package authz

import (
    `context`
    
    `github.com/gin-gonic/gin`
)

type Option func(*acl)
type RouterOption func(*gin.RouterGroup, *acl)

type AccessControl interface {
    Middleware() gin.HandlerFunc
    WithOption(...Option) AccessControl
    WithRouterOption(*gin.RouterGroup, ...RouterOption) AccessControl
}

type User interface {
    UserID() uint64
    Role() string
}

type TokenParser interface {
    Parse(ctx context.Context, token string) (User, error)
}

type PermissionChecker interface {
    Check(ctx context.Context, subject string, object string, action string) (bool, error)
}
