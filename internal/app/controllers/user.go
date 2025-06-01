package controllers

import (
    `dpcms/internal/app/errors/user_error`
    `dpcms/internal/app/middleware/auth`
    `dpcms/internal/app/services`
    `dpcms/internal/app/services/types/user`
    `dpcms/internal/erroz`
    `dpcms/internal/infra`
    "github.com/gin-gonic/gin"
)

type UserController struct {
    services *services.Services
    infra    *infra.Infra
}

func (c UserController) setup(server *gin.Engine) {
    ac, middleware := auth.New(c.services)
    
    g := server.Group("/user", middleware)
    g.POST("/register", c.Register)
    g.POST("/login", c.Login)
    g.POST("/logout", c.Logout)
    g.POST("/update-password", c.UpdatePassword)
    g.POST("/reset-password", c.ResetPassword)
    g.POST("/list", c.List)
    g.POST("/detail", c.Detail)
    g.POST("/delete", c.Delete)
    
    ac.AddWhitelist("/register", "/login")
    ac.SetObjectName("user").
        AddGroupPermission(g, "/list", "read").
        AddGroupPermission(g, "/delete", "delete").
        AddGroupPermission(g, "/update", "edit").
        AddGroupPermission(g, "/reset-password", "reset-password")
}

func NewUserController(s *services.Services, i *infra.Infra) UserController {
    return UserController{services: s, infra: i}
}

func (c UserController) Register(ctx *gin.Context) {
    params := user.RegisterParams{IP: ctx.ClientIP()}
    if err := ctx.ShouldBindJSON(&params); err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    u, err := c.services.User.Create(ctx, params)
    if err != nil {
        erroz.ResolveWithWrite(ctx, err)
    } else {
        erroz.OK.WithOption(erroz.WithData(u)).Write(ctx)
    }
}

func (c UserController) Login(ctx *gin.Context) {
    params := user.CredentialParams{}
    if err := ctx.ShouldBindJSON(&params); err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    u, valid, err := c.services.User.FindByCredential(ctx, params)
    if err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    if !valid {
        user_error.ErrWrongPassword.Write(ctx)
        return
    }
    token, err := c.services.Token.Create(u.ID)
    if err != nil {
        erroz.ResolveWithWrite(ctx, err)
    } else {
        erroz.OK.WithOption(erroz.WithData(user.LoginResult{Detail: u, Token: token})).Write(ctx)
    }
}

func (c UserController) Logout(ctx *gin.Context) {
    tokenString := ctx.GetHeader("Authorization")
    if tokenString != "" {
        _ = c.services.Token.Revoke(ctx, tokenString)
    }
    erroz.OK.Write(ctx)
}

func (c UserController) UpdatePassword(ctx *gin.Context) {
    params := user.UpdatePasswordParams{}
    if err := ctx.ShouldBindJSON(&params); err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    u := auth.GetAuthorizedUser(ctx)
    valid, err := c.services.User.IsValidCredential(ctx, u.Username, params.RawPassword)
    if err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    if !valid {
        user_error.ErrWrongPassword.Write(ctx)
        return
    }
    err = c.services.User.ResetPassword(ctx, u.ID, u.Password)
    if err != nil {
        erroz.ResolveWithWrite(ctx, err)
    } else {
        erroz.OK.Write(ctx)
    }
}

func (c UserController) ResetPassword(ctx *gin.Context) {
    params := user.ResetPasswordParams{}
    if err := ctx.ShouldBindJSON(&params); err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    userInfo, err := c.services.User.FindByID(ctx, params.UserID)
    if err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    err = c.services.User.ResetPassword(ctx, userInfo.ID, params.Password)
    if err != nil {
        erroz.ResolveWithWrite(ctx, err)
    } else {
        erroz.OK.Write(ctx)
    }
}

func (c UserController) List(ctx *gin.Context) {
    params := user.ListRetrieveParams{}
    if err := ctx.ShouldBindQuery(&params); err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    result, err := c.services.User.List(ctx, params)
    if err != nil {
        erroz.ResolveWithWrite(ctx, err)
    } else {
        erroz.OK.WithOption(erroz.WithData(result)).Write(ctx)
    }
}

func (c UserController) Detail(ctx *gin.Context) {
    u := auth.GetAuthorizedUser(ctx)
    erroz.OK.WithOption(erroz.WithData(u)).Write(ctx)
}

func (c UserController) Delete(ctx *gin.Context) {
    params := user.DeleteParams{}
    if err := ctx.ShouldBindJSON(&params); err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    err := c.services.User.Delete(ctx, params.ID)
    if err != nil {
        erroz.ResolveWithWrite(ctx, err)
    } else {
        erroz.OK.Write(ctx)
    }
}
