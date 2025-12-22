package controller

import (
    `dpcms/internal/app/controller/internal/httpbinding`
    `dpcms/internal/app/dto`
    `dpcms/internal/app/erroz`
    `dpcms/internal/app/middleware/authz`
    `dpcms/internal/app/service`
    `dpcms/internal/infra`
    `dpcms/internal/infra/logger`
    
    "github.com/gin-gonic/gin"
    `go.uber.org/zap`
)

type UserController struct {
    services *service.Services
    log      *logger.Logger
}

func NewUserController(s *service.Services, i *infra.Infra) UserController {
    return UserController{services: s, log: i.Log}
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
    httpbinding.BindJSON[dto.UserCreateParams](ctx, func(params dto.UserCreateParams) (any, error) {
        params.IP = ctx.ClientIP()
        return c.services.User.Create(ctx, params)
    })
}

func (c UserController) Login(ctx *gin.Context) {
    httpbinding.BindJSON[dto.UserCredentialParams](ctx, func(params dto.UserCredentialParams) (any, error) {
        u, err := c.services.User.FindByCredential(ctx, params)
        if err != nil {
            return nil, err
        }
        token, err := c.services.Token.Create(u.ID)
        if err != nil {
            return nil, err
        }
        return dto.UserAuthnResult{User: u, Token: token}, nil
    })
}

func (c UserController) Logout(ctx *gin.Context) {
    tokenString := ctx.GetHeader("Authorization")
    if tokenString != "" {
        err := c.services.Token.Revoke(ctx, tokenString)
        if err != nil {
            c.log.App.Warn("TOKEN注销失败", zap.Error(err))
        }
    }
    erroz.OK.Write(ctx)
}

func (c UserController) UpdatePassword(ctx *gin.Context) {
    httpbinding.BindJSON[dto.UserPasswdUpdateParams](ctx, func(params dto.UserPasswdUpdateParams) (any, error) {
        if params.Password != params.PasswordConfirm {
            return nil, erroz.UserWrongConfirmPassword.ToError()
        }
        if params.Password == params.OldPassword {
            return nil, erroz.UserNewPwdEqualsOldPwd.ToError()
        }
        u := authz.GetAuthorizedUser(ctx)
        _, err := c.services.User.FindByCredential(ctx, dto.UserCredentialParams{
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
    httpbinding.BindJSON[dto.UserPasswdResetParams](ctx, func(params dto.UserPasswdResetParams) (any, error) {
        userInfo, err := c.services.User.FindByID(ctx, params.ID)
        if err != nil {
            return nil, err
        }
        return nil, c.services.User.ResetPassword(ctx, userInfo.ID, params.Password)
    })
}

func (c UserController) List(ctx *gin.Context) {
    httpbinding.BindJSON[dto.UserListQueryParams](ctx, func(params dto.UserListQueryParams) (any, error) {
        return c.services.User.List(ctx, params)
    })
}

func (c UserController) Profile(ctx *gin.Context) {
    u := authz.GetAuthorizedUser(ctx)
    erroz.OK.WithOption(erroz.WithData(u)).Write(ctx)
}

func (c UserController) Detail(ctx *gin.Context) {
    httpbinding.BindJSON[dto.QueryByResourceID](ctx, func(params dto.QueryByResourceID) (any, error) {
        return c.services.User.FindByID(ctx, params.ID)
    })
}

func (c UserController) Delete(ctx *gin.Context) {
    httpbinding.BindJSON[dto.QueryByResourceID](ctx, func(params dto.QueryByResourceID) (any, error) {
        return nil, c.services.User.Delete(ctx, params.ID)
    })
}

func (c UserController) UpdateProfile(ctx *gin.Context) {
    httpbinding.BindJSON[dto.UserProfileUpdateParams](ctx, func(params dto.UserProfileUpdateParams) (any, error) {
        return nil, c.services.User.UpdateProfile(ctx, params)
    })
}
