package authz

import (
    `strings`
    
    `cms/internal/erroz`
    
    `github.com/armon/go-radix`
    `github.com/gin-gonic/gin`
)

type acl struct {
    token         TokenParser
    perm          PermissionChecker
    object        string
    whitelistTree *radix.Tree
    permTree      *radix.Tree
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
        if _, ok := s.whitelistTree.Get(ctx.Request.URL.Path); ok {
            ctx.Next()
            return
        }
        
        credential := ctx.GetHeader("Authorization")
        if credential == "" {
            ErrAuthorized.Abort(ctx)
            return
        }
        
        token, found := strings.CutPrefix(credential, "Bearer ")
        if !found {
            ErrAuthorized.Abort(ctx)
            return
        }
        
        user, err := s.token.Parse(ctx, token)
        if err != nil {
            erroz.ResolveWithAbort(ctx, err)
            return
        }
        setCurrentUser(ctx, user)
        
        if s.object == "" {
            ctx.Next()
            return
        }
        
        path := ctx.FullPath()
        _, perm, found := s.permTree.LongestPrefix(path)
        if found {
            permStr, ok := perm.(string)
            if !ok {
                ErrAccessDenied.Abort(ctx)
                return
            }
            
            if pass, err := s.perm.Check(ctx, user.Role(), s.object, permStr); err != nil {
                ErrAccessDenied.Abort(ctx)
                return
            } else if !pass {
                ErrAccessDenied.Abort(ctx)
                return
            }
        }
        
        ctx.Next()
    }
}
