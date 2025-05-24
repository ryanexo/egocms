package controller

import (
    "context"
    
    `dpcms/internal/database/model`
    `dpcms/internal/enum`
    erroz2 `dpcms/internal/erroz`
    `dpcms/internal/http/controller/internal/dto`
    `dpcms/internal/http/errors/user_error`
    `dpcms/internal/http/middleware/auth`
    `dpcms/internal/http/service`
    `dpcms/internal/infra`
    `dpcms/internal/packages/password`
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
        erroz2.Resolve(ctx, err)
        return
    }
    finalUser := &model.User{
        Username: *u.Username,
        Password: *u.Password,
        Email:    *u.Email,
        IP:       ctx.ClientIP(),
    }
    if err := c.service.User.Create(context.Background(), finalUser); err != nil {
        erroz2.Resolve(ctx, err)
        return
    }
    erroz2.OK.WithOption(erroz2.WithData(finalUser)).Handle(ctx)
}

func (c UserController) Grant(ctx *gin.Context) {
    u := &dto.User{}
    if err := ctx.ShouldBindJSON(u); err != nil {
        erroz2.Resolve(ctx, err)
        return
    }
    result, err := c.service.User.FindByName(context.Background(), *u.Username)
    if err != nil {
        erroz2.Resolve(ctx, err)
        return
    }
    if !password.Compare(result.Password, *u.Password) {
        user_error.ErrWrongPassword.Handle(ctx)
        return
    }
    t, err := c.service.Token.Create(result.ID)
    if err != nil {
        erroz2.Resolve(ctx, err)
        return
    }
    erroz2.OK.WithOption(erroz2.WithData(map[string]string{"token": t})).Handle(ctx)
}

func (c UserController) Revoke(ctx *gin.Context) {
    u := auth.GetAuthorizedUser(ctx)
    key := enum.GetCacheUserBlacklistKey(u.ID)
    c.infra.Cache.Set(key, true)
    erroz2.OK.Handle(ctx)
}

func (c UserController) UpdatePassword(ctx *gin.Context) {
    u := &dto.User{}
    if err := ctx.ShouldBindJSON(u); err != nil {
        erroz2.Resolve(ctx, err)
        return
    }
    userInfo := auth.GetAuthorizedUser(ctx)
    err := c.service.User.UpdatePassword(context.Background(), userInfo.ID, *u.Password)
    if err != nil {
        erroz2.Resolve(ctx, err)
        return
    }
    erroz2.OK.Handle(ctx)
}
