package auth

import (
    `context`
    `strings`
    
    `dpcms/api/service`
    `dpcms/enum`
    `dpcms/erroz`
    `dpcms/model`
    "github.com/gin-gonic/gin"
)

type Auth struct {
    services  *service.Services
    whitelist map[string]bool
}

func New(srv *service.Services) *Auth {
    return &Auth{services: srv, whitelist: make(map[string]bool)}
}

func (auth *Auth) CreateMiddleware() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        if auth.whitelist[ctx.Request.URL.Path] {
            ctx.Next()
            return
        }
        
        credential := ctx.GetHeader("Authorization")
        if credential == "" {
            erroz.ErrUnauthorized.Abort(ctx)
            return
        }
        
        token, found := strings.CutPrefix(credential, "Bearer ")
        if !found {
            erroz.ErrUnauthorized.Abort(ctx)
            return
        }
        
        err := auth.shouldSetUserWithToken(ctx, token)
        if err != nil {
            erroz.ResolveWithAbort(ctx, err)
            return
        }
        ctx.Next()
    }
}

func (auth *Auth) shouldSetUserWithToken(ctx *gin.Context, token string) error {
    user, err := auth.services.Token.Parse(context.Background(), token)
    if err != nil {
        return err
    }
    ctx.Set(enum.ApiAuthCurrentUser, user)
    return nil
}

func (auth *Auth) Ignore(path ...string) *Auth {
    for _, p := range path {
        auth.whitelist[p] = true
    }
    return auth
}

func (auth *Auth) IgnoreWithPrefix(prefix string, path ...string) *Auth {
    for _, p := range path {
        currentPath := prefix + p
        auth.whitelist[currentPath] = true
    }
    return auth
}

func GetAuthorizedUser(ctx *gin.Context) *model.User {
    return ctx.MustGet(enum.ApiAuthCurrentUser).(*model.User)
}
