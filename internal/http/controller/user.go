package controller

import (
    "context"
    
    `dpcms/internal/enum`
    `dpcms/internal/erroz`
    `dpcms/internal/http/dto`
    `dpcms/internal/http/middleware/auth`
    `dpcms/internal/http/service`
    `dpcms/internal/infra`
    "github.com/gin-gonic/gin"
)

type UserController struct {
    service *service.Services
    infra   *infra.Infra
}

func (c UserController) setup(server *gin.Engine) {
    g := server.Group("/user")
    g.POST("/register", c.Register)
    g.POST("/login", c.Login)
    g.POST("/logout", c.Logout)
    g.POST("/update-password", c.UpdatePassword)
    
    auth.New(c.service).
        SkipWithGroup(g, []string{"/register", "/login"}).
        Append(g)
}

func NewUserController(s *service.Services, i *infra.Infra) *UserController {
    return &UserController{service: s, infra: i}
}

func (c UserController) Register(ctx *gin.Context) {
    ctx.FullPath()
    u := &dto.UserRegisterRequest{IP: ctx.ClientIP()}
    if err := ctx.ShouldBindJSON(u); err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    userInfo, err := c.service.User.Create(context.Background(), u)
    if err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    erroz.OK.WithOption(erroz.WithData(userInfo)).Write(ctx)
}

func (c UserController) Login(ctx *gin.Context) {
    u := &dto.UserLoginRequest{}
    if err := ctx.ShouldBindJSON(u); err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    userInfo, err := c.service.User.FindUserWithCredential(context.Background(), *u.Username, *u.Password)
    if err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    token, err := c.service.Token.Create(userInfo.ID)
    if err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    erroz.OK.WithOption(erroz.WithData(dto.UserLoginResponse{Token: token})).Write(ctx)
}

func (c UserController) Logout(ctx *gin.Context) {
    u := auth.GetAuthorizedUser(ctx)
    c.service.Token.Revoke()
    erroz.OK.Write(ctx)
}

func (c UserController) UpdatePassword(ctx *gin.Context) {
    u := &dto.UserRegisterRequest{}
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
