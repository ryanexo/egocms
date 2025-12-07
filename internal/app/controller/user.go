package controller

import (
    erroz2 `dpcms/internal/app/erroz`
    `dpcms/internal/app/middleware/authz`
    `dpcms/internal/app/service`
    `dpcms/internal/app/service/types/user`
    `dpcms/internal/infra`
    
    "github.com/gin-gonic/gin"
)

type UserController struct {
    services *service.Services
    infra    *infra.Infra
}

func (c UserController) setup(server *gin.Engine) {
    acl := authz.New(c.services)
    
    g := server.Group("/user", acl.Middleware())
    g.POST("/register", c.Register)
    g.POST("/login", c.Login)
    g.POST("/logout", c.Logout)
    g.POST("/update-password", c.UpdatePassword)
    g.POST("/reset-password", c.ResetPassword)
    g.POST("/list", c.List)
    g.POST("/detail", c.Detail)
    g.POST("/delete", c.Delete)
    g.POST("/update-profile", c.UpdateProfile)
    
    acl.WithRouterOption(g, authz.WithRouterWhitelist("/register", "/login"))
    acl.WithRouterOption(g,
        authz.WithRouterPermission("/list", "read"),
        authz.WithRouterPermission("/delete", "delete"),
        authz.WithRouterPermission("/update-profile", "update"),
        authz.WithRouterPermission("/update-password", "update"),
        authz.WithRouterPermission("/reset-password", "reset-password"),
    )
}

func NewUserController(s *service.Services, i *infra.Infra) UserController {
    return UserController{services: s, infra: i}
}

func (c UserController) Register(ctx *gin.Context) {
    params := user.CreateDTO{IP: ctx.ClientIP()}
    if err := ctx.ShouldBindJSON(&params); err != nil {
        erroz2.ResolveWithWrite(ctx, err)
        return
    }
    u, err := c.services.User.Create(ctx, params)
    if err != nil {
        erroz2.ResolveWithWrite(ctx, err)
    } else {
        erroz2.OK.WithOption(erroz2.WithData(u)).Write(ctx)
    }
}

func (c UserController) Login(ctx *gin.Context) {
    params := user.CredentialDTO{}
    if err := ctx.ShouldBindJSON(&params); err != nil {
        erroz2.ResolveWithWrite(ctx, err)
        return
    }
    u, valid, err := c.services.User.FindByCredential(ctx, params)
    if err != nil {
        erroz2.ResolveWithWrite(ctx, err)
        return
    }
    if !valid {
        erroz2.ErrWrongPassword.Write(ctx)
        return
    }
    token, err := c.services.Token.Create(u.ID)
    if err != nil {
        erroz2.ResolveWithWrite(ctx, err)
    } else {
        erroz2.OK.WithOption(erroz2.WithData(user.LoginResult{Detail: u, Token: token})).Write(ctx)
    }
}

func (c UserController) Logout(ctx *gin.Context) {
    tokenString := ctx.GetHeader("Authorization")
    if tokenString != "" {
        _ = c.services.Token.Revoke(ctx, tokenString)
    }
    erroz2.OK.Write(ctx)
}

func (c UserController) UpdatePassword(ctx *gin.Context) {
    params := user.UpdatePasswordDTO{}
    if err := ctx.ShouldBindJSON(&params); err != nil {
        erroz2.ResolveWithWrite(ctx, err)
        return
    }
    u := authz.GetAuthorizedUser(ctx)
    valid, err := c.services.User.IsValidCredential(ctx, u.Username, params.RawPassword)
    if err != nil {
        erroz2.ResolveWithWrite(ctx, err)
        return
    }
    if !valid {
        erroz2.ErrWrongPassword.Write(ctx)
        return
    }
    err = c.services.User.ResetPassword(ctx, u.ID, u.Password)
    if err != nil {
        erroz2.ResolveWithWrite(ctx, err)
    } else {
        erroz2.OK.Write(ctx)
    }
}

func (c UserController) ResetPassword(ctx *gin.Context) {
    params := user.ResetPasswordDTO{}
    if err := ctx.ShouldBindJSON(&params); err != nil {
        erroz2.ResolveWithWrite(ctx, err)
        return
    }
    userInfo, err := c.services.User.FindByID(ctx, params.ID)
    if err != nil {
        erroz2.ResolveWithWrite(ctx, err)
        return
    }
    err = c.services.User.ResetPassword(ctx, userInfo.ID, params.Password)
    if err != nil {
        erroz2.ResolveWithWrite(ctx, err)
    } else {
        erroz2.OK.Write(ctx)
    }
}

func (c UserController) List(ctx *gin.Context) {
    params := user.ListDTO{}
    if err := ctx.ShouldBindQuery(&params); err != nil {
        erroz2.ResolveWithWrite(ctx, err)
        return
    }
    result, err := c.services.User.List(ctx, params)
    if err != nil {
        erroz2.ResolveWithWrite(ctx, err)
    } else {
        erroz2.OK.WithOption(erroz2.WithData(result)).Write(ctx)
    }
}

func (c UserController) Detail(ctx *gin.Context) {
    u := authz.GetAuthorizedUser(ctx)
    erroz2.OK.WithOption(erroz2.WithData(u)).Write(ctx)
}

func (c UserController) Delete(ctx *gin.Context) {
    params := user.DeleteDTO{}
    if err := ctx.ShouldBindJSON(&params); err != nil {
        erroz2.ResolveWithWrite(ctx, err)
        return
    }
    err := c.services.User.Delete(ctx, params.ID)
    if err != nil {
        erroz2.ResolveWithWrite(ctx, err)
    } else {
        erroz2.OK.Write(ctx)
    }
}

func (c UserController) UpdateProfile(ctx *gin.Context) {
    params := user.UpdateProfileDTO{}
    if err := ctx.ShouldBindJSON(&params); err != nil {
        erroz2.ResolveWithWrite(ctx, err)
        return
    }
    err := c.services.User.UpdateProfile(ctx, params)
    if err != nil {
        erroz2.ResolveWithWrite(ctx, err)
    } else {
        erroz2.OK.Write(ctx)
    }
}
