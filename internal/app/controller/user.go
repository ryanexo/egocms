package controller

import (
    `dpcms/internal/app/controller/internal/common`
    `dpcms/internal/app/erroz`
    `dpcms/internal/app/middleware/authz`
    `dpcms/internal/app/service`
    `dpcms/internal/app/service/srvparams`
    `dpcms/internal/infra`
    
    "github.com/gin-gonic/gin"
    `go.uber.org/zap`
)

type UserController struct {
    services *service.Services
    infra    *infra.Infra
}

func NewUserController(s *service.Services, i *infra.Infra) UserController {
    return UserController{services: s, infra: i}
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
    g.POST("/profile", c.Profile)
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

func (c UserController) Register(ctx *gin.Context) {
    common.BindJSON[srvparams.UserCreateParams](ctx, func(params srvparams.UserCreateParams) (any, error) {
        params.IP = ctx.ClientIP()
        return c.services.User.Create(ctx, params)
    })
}

func (c UserController) Login(ctx *gin.Context) {
    common.BindJSON[srvparams.UserCredentialParams](ctx, func(params srvparams.UserCredentialParams) (any, error) {
        u, err := c.services.User.FindByCredential(ctx, params)
        if err != nil {
            return nil, err
        }
        token, err := c.services.Token.Create(u.ID)
        if err != nil {
            return nil, err
        }
        return srvparams.UserAuthnResult{User: u, Token: token}, nil
    })
}

func (c UserController) Logout(ctx *gin.Context) {
    tokenString := ctx.GetHeader("Authorization")
    if tokenString != "" {
        err := c.services.Token.Revoke(ctx, tokenString)
        if err != nil {
            c.infra.Log.App.Warn("TOKEN注销失败", zap.Error(err))
        }
    }
    erroz.OK.Write(ctx)
}

func (c UserController) UpdatePassword(ctx *gin.Context) {
    common.BindJSON[srvparams.UserPasswdUpdateParams](ctx, func(params srvparams.UserPasswdUpdateParams) (any, error) {
        if params.Password != params.PasswordConfirm {
            return nil, erroz.ErrWrongConfirmPassword.ToError()
        }
        if params.Password == params.OldPassword {
            return nil, erroz.ErrNewPwdEqualsOldPwd.ToError()
        }
        u := authz.GetAuthorizedUser(ctx)
        _, err := c.services.User.FindByCredential(ctx, srvparams.UserCredentialParams{
            Username: u.Username,
            Password: params.OldPassword,
        })
        if err != nil {
            return nil, err
        }
        return nil, c.services.User.ResetPassword(ctx, u.ID, u.Password)
    })
}

func (c UserController) ResetPassword(ctx *gin.Context) {
    common.BindJSON[srvparams.UserPasswdResetParams](ctx, func(params srvparams.UserPasswdResetParams) (any, error) {
        userInfo, err := c.services.User.FindByID(ctx, params.ID)
        if err != nil {
            return nil, err
        }
        return nil, c.services.User.ResetPassword(ctx, userInfo.ID, params.Password)
    })
}

func (c UserController) List(ctx *gin.Context) {
    common.BindJSON[srvparams.UserListQueryParams](ctx, func(params srvparams.UserListQueryParams) (any, error) {
        return c.services.User.List(ctx, params)
    })
}

func (c UserController) Profile(ctx *gin.Context) {
    u := authz.GetAuthorizedUser(ctx)
    erroz.OK.WithOption(erroz.WithData(u)).Write(ctx)
}

func (c UserController) Detail(ctx *gin.Context) {
    common.BindJSON[srvparams.QueryByResourceID](ctx, func(params srvparams.QueryByResourceID) (any, error) {
        return c.services.User.FindByID(ctx, params.ID)
    })
}

func (c UserController) Delete(ctx *gin.Context) {
    common.BindJSON[srvparams.QueryByResourceID](ctx, func(params srvparams.QueryByResourceID) (any, error) {
        return nil, c.services.User.Delete(ctx, params.ID)
    })
}

func (c UserController) UpdateProfile(ctx *gin.Context) {
    common.BindJSON[srvparams.UserProfileUpdateParams](ctx, func(params srvparams.UserProfileUpdateParams) (any, error) {
        return nil, c.services.User.UpdateProfile(ctx, params)
    })
}
