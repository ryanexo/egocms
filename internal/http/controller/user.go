package controller

import (
    "context"
    
    `dpcms/internal/database/model`
    `dpcms/internal/enum`
    `dpcms/internal/erroz`
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
    middleware := auth.New(c.service).SkipWithPrefix("/user", []string{"/register", "/login"})
    
    g := server.Group("/user", middleware.CreateMiddleware())
    g.POST("/register", c.Register)
    g.POST("/grant", c.Login)
    g.POST("/revoke", c.Logout)
    g.POST("/update-password", c.UpdatePassword)
}

func NewUserController(s *service.Services, i *infra.Infra) *UserController {
    return &UserController{service: s, infra: i}
}

func (c UserController) Register(ctx *gin.Context) {
    ctx.FullPath()
    u := &dto.UserAuth{}
    if err := ctx.ShouldBindJSON(u); err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    finalUser := &model.User{
        Username: *u.Username,
        Password: *u.Password,
        Email:    *u.Email,
        IP:       ctx.ClientIP(),
    }
    if err := c.service.User.Create(context.Background(), finalUser); err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    erroz.OK.WithOption(erroz.WithData(finalUser)).Write(ctx)
}

func (c UserController) Login(ctx *gin.Context) {
    u := &dto.UserAuth{}
    u.SetValidationFields([]string{"username", "password"})
    if err := ctx.ShouldBindJSON(u); err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    result, err := c.service.User.FindByName(context.Background(), *u.Username)
    if err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    if !password.Compare(result.Password, *u.Password) {
        user_error.ErrWrongPassword.Write(ctx)
        return
    }
    t, err := c.service.Token.Create(result.ID)
    if err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    erroz.OK.WithOption(erroz.WithData(map[string]string{"token": t})).Write(ctx)
}

func (c UserController) Logout(ctx *gin.Context) {
    u := auth.GetAuthorizedUser(ctx)
    key := enum.GetCacheUserBlacklistKey(u.ID)
    c.infra.Cache.Set(key, true)
    erroz.OK.Write(ctx)
}

func (c UserController) UpdatePassword(ctx *gin.Context) {
    u := &dto.UserAuth{}
    if err := ctx.ShouldBindJSON(u); err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    userInfo := auth.GetAuthorizedUser(ctx)
    err := c.service.User.UpdatePassword(context.Background(), userInfo.ID, *u.Password)
    if err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    erroz.OK.Write(ctx)
}
