package controller

import (
    `dpcms/internal/app/controller/internal/common`
    `dpcms/internal/app/erroz`
    `dpcms/internal/app/middleware/authz`
    `dpcms/internal/app/service`
    `dpcms/internal/app/service/srvparams`
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
    common.BasicBind[srvparams.UserCreateParams](ctx, func(params srvparams.UserCreateParams) (any, error) {
        params.IP = ctx.ClientIP()
        return c.services.User.Create(ctx, params)
    })
}

func (c UserController) Login(ctx *gin.Context) {
    common.BasicBind[srvparams.UserCredentialParams](ctx, func(params srvparams.UserCredentialParams) (any, error) {
        u, valid, err := c.services.User.FindByCredential(ctx, params)
    })
    P := srvparams.UserCredentialParams{}
    if err := ctx.ShouldBindJSON(&P); err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    
    if err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    if !valid {
        erroz.ErrWrongPassword.Write(ctx)
        return
    }
    token, err := c.services.Token.Create(u.ID)
    if err != nil {
        erroz.ResolveWithWrite(ctx, err)
    } else {
        erroz.OK.WithOption(erroz.WithData(srvparams.UserAuthnResult{User: u, Token: token})).Write(ctx)
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
    p := srvparams.UserPasswdUpdateParams{}
    if err := ctx.ShouldBindJSON(&p); err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    u := authz.GetAuthorizedUser(ctx)
    valid, err := c.services.User.FindByCredential(ctx, u.Username, p.RawPassword)
    if err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    if !valid {
        erroz.ErrWrongPassword.Write(ctx)
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
    p := srvparams.UserPasswdResetParams{}
    if err := ctx.ShouldBindJSON(&p); err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    userInfo, err := c.services.User.FindByID(ctx, p.ID)
    if err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    err = c.services.User.ResetPassword(ctx, userInfo.ID, p.Password)
    if err != nil {
        erroz.ResolveWithWrite(ctx, err)
    } else {
        erroz.OK.Write(ctx)
    }
}

func (c UserController) List(ctx *gin.Context) {
    p := srvparams.UserListQueryParams{}
    if err := ctx.ShouldBindQuery(&p); err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    result, err := c.services.User.List(ctx, p)
    if err != nil {
        erroz.ResolveWithWrite(ctx, err)
    } else {
        erroz.OK.WithOption(erroz.WithData(result)).Write(ctx)
    }
}

func (c UserController) Detail(ctx *gin.Context) {
    u := authz.GetAuthorizedUser(ctx)
    erroz.OK.WithOption(erroz.WithData(u)).Write(ctx)
}

func (c UserController) Delete(ctx *gin.Context) {
    p := srvparams.QueryByResourceID{}
    if err := ctx.ShouldBindJSON(&p); err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    err := c.services.User.Delete(ctx, p.ID)
    if err != nil {
        erroz.ResolveWithWrite(ctx, err)
    } else {
        erroz.OK.Write(ctx)
    }
}

func (c UserController) UpdateProfile(ctx *gin.Context) {
    p := srvparams.UserProfileUpdateParams{}
    if err := ctx.ShouldBindJSON(&p); err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    err := c.services.User.UpdateProfile(ctx, p)
    if err != nil {
        erroz.ResolveWithWrite(ctx, err)
    } else {
        erroz.OK.Write(ctx)
    }
}
