package controller

import (
    "context"
    
    `dpcms/api/controller/internal/dto`
    `dpcms/api/infra`
    `dpcms/api/middleware/auth`
    `dpcms/api/service`
    `dpcms/config`
    `dpcms/enum`
    `dpcms/erroz`
    `dpcms/model`
    `dpcms/packages/password`
    "github.com/gin-gonic/gin"
)

type UserController struct {
    service *service.Services
    infra   *infra.Infra
}

func (c UserController) setup(server *gin.Engine) {
    authMiddleware := auth.New(c.service).IgnoreWithPrefix("/user", "/register", "/login").CreateMiddleware()
    g := server.Group("/user", authMiddleware)
    g.POST("/register", c.Register)
    g.POST("/grant", c.Grant)
    g.POST("/revoke", c.Revoke)
    g.POST("/update-password", c.UpdatePassword)
}

func NewUserController(s *service.Services, i *infra.Infra) *UserController {
    return &UserController{service: s, infra: i}
}

func (c UserController) Register(ctx *gin.Context) {
    u := &dto.User{}
    if err := ctx.ShouldBindJSON(u); err != nil {
        erroz.Resolve(ctx, err)
        return
    }
    finalUser := &model.User{
        Username: *u.Username,
        Password: *u.Password,
        Email:    *u.Email,
        IP:       ctx.ClientIP(),
    }
    if err := c.service.User.Create(context.Background(), finalUser); err != nil {
        erroz.Resolve(ctx, err)
        return
    }
    erroz.OK.WithOption(erroz.WithData(finalUser)).Apply(ctx)
}

func (c UserController) Grant(ctx *gin.Context) {
    u := &dto.User{}
    if err := ctx.ShouldBindJSON(u); err != nil {
        erroz.Resolve(ctx, err)
        return
    }
    result, err := c.service.User.FindByName(context.Background(), *u.Username)
    if err != nil {
        erroz.Resolve(ctx, err)
        return
    }
    if !password.Compare(result.Password, *u.Password) {
        erroz.ErrWrongPassword.Apply(ctx)
        return
    }
    claims := &dto.UserToken{UserID: result.ID}
    t, err := claims.Create(config.Get().GlobalKey)
    if err != nil {
        erroz.Resolve(ctx, err)
        return
    }
    erroz.OK.WithOption(erroz.WithData(map[string]string{"token": t})).Apply(ctx)
}

func (c UserController) Revoke(ctx *gin.Context) {
    u := auth.GetAuthorizedUser(ctx)
    key := enum.GetCacheUserBlacklistKey(u.ID)
    c.infra.Cache.Set(key, true)
    erroz.OK.Apply(ctx)
}

func (c UserController) UpdatePassword(ctx *gin.Context) {
    u := &dto.User{}
    if err := ctx.ShouldBindJSON(u); err != nil {
        erroz.Resolve(ctx, err)
        return
    }
    userInfo := auth.GetAuthorizedUser(ctx)
    err := c.service.User.UpdatePassword(context.Background(), userInfo.ID, *u.Password)
    if err != nil {
        erroz.Resolve(ctx, err)
        return
    }
    erroz.OK.Apply(ctx)
}
