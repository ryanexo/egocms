package auth

import (
    `context`
    `strings`
    
    `dpcms/api/dto`
    `dpcms/api/service`
    `dpcms/config`
    `dpcms/enum`
    `dpcms/erroz`
    `dpcms/model`
    "github.com/gin-gonic/gin"
    `github.com/golang-jwt/jwt/v5`
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
    userClaims := &dto.UserToken{}
    cfg := config.Get()
    tokenResult, _ := jwt.ParseWithClaims(token, userClaims, func(t *jwt.Token) (interface{}, error) {
        return []byte(cfg.GlobalKey), nil
    })
    if !tokenResult.Valid {
        return erroz.ErrUnauthorized.ToError()
    }
    user, err := auth.services.User.FindByID(context.Background(), userClaims.UserID)
    if err != nil {
        return err
    }
    ctx.Set(enum.CurrentUser, user)
    return nil
}

func (auth *Auth) Ignore(path ...string) *Auth {
    for _, p := range path {
        auth.whitelist[p] = true
    }
    return auth
}

func GetAuthorizedUser(ctx *gin.Context) *model.User {
    return ctx.MustGet(enum.CurrentUser).(*model.User)
}
