package auth

import (
    `context`
    `strings`
    
    `dpcms/internal/database/model`
    `dpcms/internal/enum`
    `dpcms/internal/erroz`
    `dpcms/internal/http/errors/auth_error`
    `dpcms/internal/http/helper/rbachelper`
    `dpcms/internal/http/service`
    "github.com/gin-gonic/gin"
)

type Auth struct {
    name       string
    services   *service.Services
    whitelist  *whitelist
    permission map[string]string
}

type RouteResource interface {
    Use(...gin.HandlerFunc) gin.IRoutes
}

func New(srv *service.Services) *Auth {
    return &Auth{services: srv, whitelist: newWhitelist()}
}

func (auth *Auth) createMiddleware() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        if auth.whitelist.match(ctx.Request.URL.Path) {
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
        
        user, err := auth.shouldSetUserWithToken(ctx, token)
        if err != nil {
            erroz.ResolveWithAbort(ctx, err)
            return
        }
        
        if auth.name != "" {
            roleName := rbachelper.GetRoleSubject(user.RoleID)
            path := ctx.FullPath()
            permission, found := auth.permission[path]
            if found {
                pass, err := auth.services.RBAC.Enforce(roleName, auth.name, permission)
                if err != nil {
                    erroz.ResolveWithAbort(ctx, err)
                    return
                }
                if !pass {
                    auth_error.ErrUnauthorized.WriteWithAbort(ctx)
                    return
                }
            }
        }
        
        ctx.Next()
    }
}

func (auth *Auth) Append(route RouteResource) {
    middleware := auth.createMiddleware()
    route.Use(middleware)
}

func (auth *Auth) SetSourceName(name string) *Auth {
    auth.name = name
    return auth
}

func (auth *Auth) Skip(path []string) *Auth {
    for _, p := range path {
        auth.whitelist.insert(p)
    }
    return auth
}

func (auth *Auth) SkipWithGroup(group *gin.RouterGroup, path []string) *Auth {
    for _, p := range path {
        mergedPath := group.BasePath() + p
        auth.whitelist.insert(mergedPath)
    }
    return auth
}

func (auth *Auth) Permission(permission map[string]string) *Auth {
    auth.permission = permission
    return auth
}

func (auth *Auth) shouldSetUserWithToken(ctx *gin.Context, token string) (*model.User, error) {
    user, err := auth.services.Token.Parse(context.Background(), token)
    if err != nil {
        return nil, err
    }
    ctx.Set(enum.ApiAuthCurrentUser, user)
    return user, nil
}

func GetAuthorizedUser(ctx *gin.Context) *model.User {
    return ctx.MustGet(enum.ApiAuthCurrentUser).(*model.User)
}
