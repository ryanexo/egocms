package auth

import (
    `strings`
    
    `dpcms/internal/database/model`
    `dpcms/internal/enum`
    `dpcms/internal/erroz`
    `dpcms/internal/http/errors/auth_error`
    `dpcms/internal/http/helper/rbachelper`
    `dpcms/internal/http/middleware/auth/internal/trie`
    `dpcms/internal/http/service`
    "github.com/gin-gonic/gin"
)

type AccessControl interface {
    SetObjectName(object string) AccessControl
    AddPermission(obj string, perm string) AccessControl
    AddPermissions(map[string]string) AccessControl
    AddWhitelist(string) AccessControl
    AddWhitelists([]string) AccessControl
}

type Auth interface {
    RouterGroup(group *gin.RouterGroup) AccessControl
    Routes(gin.IRoutes, ...gin.IRoutes) AccessControl
}

type accessControl struct {
    objectName  string
    services    *service.Services
    whitelist   *trie.Trie
    permissions map[string]string
    group       *gin.RouterGroup
    routes      []gin.IRoutes
}

func New(srv *service.Services) Auth {
    return &accessControl{services: srv, whitelist: trie.NewPathTrie(), permissions: make(map[string]string), routes: make([]gin.IRoutes, 0)}
}

func (ac *accessControl) RouterGroup(group *gin.RouterGroup) AccessControl {
    ac.group = group
    middleware := ac.createMiddleware()
    group.Use(middleware)
    return ac
}

func (ac *accessControl) Routes(route gin.IRoutes, routes ...gin.IRoutes) AccessControl {
    ac.routes = append(ac.routes, route)
    if len(routes) > 0 {
        ac.routes = append(ac.routes, routes...)
    }
    middleware := ac.createMiddleware()
    for _, r := range ac.routes {
        r.Use(middleware)
    }
    return ac
}

func (ac *accessControl) SetObjectName(objectName string) AccessControl {
    ac.objectName = objectName
    return ac
}

func (ac *accessControl) AddPermission(obj string, perm string) AccessControl {
    if ac.group != nil {
        perm = ac.group.BasePath() + perm
    }
    ac.permissions[obj] = perm
    return ac
}

func (ac *accessControl) AddPermissions(permissions map[string]string) AccessControl {
    for path, perm := range permissions {
        ac.AddPermission(path, perm)
    }
    return ac
}

func (ac *accessControl) AddWhitelist(path string) AccessControl {
    if ac.group != nil {
        path = ac.group.BasePath() + path
    }
    ac.whitelist.Insert(path)
    return ac
}

func (ac *accessControl) AddWhitelists(pathList []string) AccessControl {
    for _, path := range pathList {
        ac.AddWhitelist(path)
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
        
        if ac.objectName != "" {
            subject := rbachelper.GetRoleSubject(user.RoleID)
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

func shouldSetUserFromToken(ctx *gin.Context, tokenSrv *service.TokenService, token string) (*model.User, error) {
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
