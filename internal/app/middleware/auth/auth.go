package auth

import (
    `strings`
    
    `dpcms/internal/database/model`
    `dpcms/internal/enum`
    `dpcms/internal/erroz`
    `dpcms/internal/app/errors/auth_error`
    `dpcms/internal/app/helpers/rbac`
    `dpcms/internal/app/middleware/auth/internal/trie`
    `dpcms/internal/app/services`
    "github.com/gin-gonic/gin"
)

type AccessControl interface {
    SetObjectName(object string) ObjectAccessControl
    AddWhitelist(...string) AccessControl
    AddGroupWhitelist(*gin.RouterGroup, ...string) AccessControl
}

type ObjectAccessControl interface {
    AddPermission(path string, perm string) ObjectAccessControl
    AddGroupPermission(group *gin.RouterGroup, path string, perm string) ObjectAccessControl
}

type accessControl struct {
    objectName  string
    services    *services.Services
    whitelist   *trie.Trie
    permissions map[string]string
}

func New(srv *services.Services) (AccessControl, gin.HandlerFunc) {
    ac := &accessControl{services: srv, whitelist: trie.NewPathTrie(), permissions: make(map[string]string)}
    return ac, ac.createMiddleware()
}

func (ac *accessControl) SetObjectName(objectName string) ObjectAccessControl {
    ac.objectName = objectName
    return ac
}

func (ac *accessControl) AddPermission(path string, perm string) ObjectAccessControl {
    ac.permissions[path] = perm
    return ac
}

func (ac *accessControl) AddGroupPermission(group *gin.RouterGroup, path string, perm string) ObjectAccessControl {
    p := group.BasePath() + path
    ac.AddPermission(p, perm)
    return ac
}

func (ac *accessControl) AddWhitelist(path ...string) AccessControl {
    for _, p := range path {
        ac.whitelist.Insert(p)
    }
    return ac
}

func (ac *accessControl) AddGroupWhitelist(group *gin.RouterGroup, path ...string) AccessControl {
    basePath := group.BasePath()
    for _, p := range path {
        ac.AddWhitelist(basePath + p)
    }
    return ac
}

func (ac *accessControl) createMiddleware() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        if ac.whitelist.Match(ctx.Request.URL.Path) {
            ctx.Next()
            return
        }
        
        credential := ctx.GetHeader("Authorization")
        if credential == "" {
            auth_error.ErrUnauthorized.WriteWithAbort(ctx)
            return
        }
        
        token, found := strings.CutPrefix(credential, "Bearer ")
        if !found {
            auth_error.ErrUnauthorized.WriteWithAbort(ctx)
            return
        }
        
        user, err := shouldSetUserFromToken(ctx, ac.services.Token, token)
        if err != nil {
            erroz.ResolveWithAbort(ctx, err)
            return
        }
        
        if user.ID != 1 && ac.objectName != "" {
            subject := rbac.GetRoleSubject(user.RoleID)
            path := ctx.FullPath()
            if perm, found := ac.permissions[path]; found {
                if pass, err := ac.services.RBAC.Enforce(subject, ac.objectName, perm); err != nil {
                    erroz.ResolveWithAbort(ctx, err)
                    return
                } else if !pass {
                    auth_error.ErrUnauthorized.WriteWithAbort(ctx)
                    return
                }
            }
        }
        
        ctx.Next()
    }
}

func shouldSetUserFromToken(ctx *gin.Context, tokenSrv services.TokenService, token string) (*model.User, error) {
    user, err := tokenSrv.GetUserFromToken(ctx, token)
    if err != nil {
        return nil, err
    }
    ctx.Set(enum.ApiAuthCurrentUser, user)
    return user, nil
}

func GetAuthorizedUser(ctx *gin.Context) *model.User {
    return ctx.MustGet(enum.ApiAuthCurrentUser).(*model.User)
}
