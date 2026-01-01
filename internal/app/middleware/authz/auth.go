package authz

import (
    `strings`
    
    `dpcms/internal/app/constant`
    `dpcms/internal/app/erroz`
    `dpcms/internal/app/service`
    `dpcms/internal/app/util/rbacutil`
    `dpcms/internal/infra/persistence/model`
    
    `github.com/armon/go-radix`
    "github.com/gin-gonic/gin"
)

type Option func(*acl)
type RouterOption func(*gin.RouterGroup, *acl)

type AccessControl interface {
    Middleware() gin.HandlerFunc
    WithOption(...Option) AccessControl
    WithRouterOption(*gin.RouterGroup, ...RouterOption) AccessControl
}

type acl struct {
    object    string
    service   *service.Services
    whitelist *radix.Tree
    perm      *radix.Tree
}

func New(srv *service.Services) AccessControl {
    return acl{service: srv, whitelist: radix.New(), perm: radix.New()}
}

func NewWithRBAC(srv *service.Services, object string) AccessControl {
    instance := New(srv).(acl)
    instance.object = object
    return instance
}

func (s acl) WithOption(opts ...Option) AccessControl {
    for _, opt := range opts {
        opt(&s)
    }
    return s
}

func (s acl) WithRouterOption(rg *gin.RouterGroup, opts ...RouterOption) AccessControl {
    for _, opt := range opts {
        opt(rg, &s)
    }
    return s
}

func (s acl) Middleware() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        if _, ok := s.whitelist.Get(ctx.Request.URL.Path); ok {
            ctx.Next()
            return
        }
        
        credential := ctx.GetHeader("Authorization")
        if credential == "" {
            erroz.Unauthorized.WriteWithAbort(ctx)
            return
        }
        
        token, found := strings.CutPrefix(credential, "Bearer ")
        if !found {
            erroz.Unauthorized.WriteWithAbort(ctx)
            return
        }
        
        user, err := shouldSetUserFromToken(ctx, s.service.Token, token)
        if err != nil {
            erroz.ResolveWithAbort(ctx, err)
            return
        }
        
        if user.ID != 1 && s.object != "" {
            subject := rbacutil.GetRoleSubject(user.RoleID.Raw())
            path := ctx.FullPath()
            if _, perm, found := s.perm.LongestPrefix(path); found {
                if pass, err := s.service.RBAC.GetEnforcer().Enforce(subject, s.object, perm); err != nil {
                    erroz.ResolveWithAbort(ctx, err)
                    return
                } else if !pass {
                    erroz.Unauthorized.WriteWithAbort(ctx)
                    return
                }
            }
        }
        
        ctx.Next()
    }
}

func shouldSetUserFromToken(ctx *gin.Context, tokenSrv *service.Token, token string) (*model.User, error) {
    user, err := tokenSrv.GetUserFromToken(ctx, token)
    if err != nil {
        return nil, err
    }
    ctx.Set(constant.RequestUserKey, user)
    return user, nil
}

func GetAuthorizedUser(ctx *gin.Context) *model.User {
    return ctx.MustGet(constant.RequestUserKey).(*model.User)
}
