package controller

import (
    token `dpcms/internal/app/token/service`
    `dpcms/internal/app/user/internal/dto`
    user `dpcms/internal/app/user/service`
    "dpcms/internal/erroz"
    "dpcms/internal/infra"
    `dpcms/internal/middleware/authz`
    `dpcms/internal/types`
    `dpcms/internal/util/contextutil`
    `dpcms/internal/util/httpbinding`
    
    _ `dpcms/internal/app/user/swagger`
    
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type UserController struct {
    UserSrv  *user.UserService
    TokenSrv *token.TokenService
    Infra    *infra.Infra
    Auth     *authz.Builder
}

func (s UserController) setup(server *gin.Engine) {
    acl := s.Auth.AccessControl("user")
    
    g := server.Group("/user", acl.Middleware())
    g.POST("/register", s.Register)
    g.POST("/login", s.Login)
    g.POST("/logout", s.Logout)
    g.POST("/update-password", s.ChangePassword)
    g.POST("/reset-password", s.ResetPassword)
    g.POST("/list", s.List)
    g.POST("/detail", s.Detail)
    g.POST("/delete", s.Delete)
    g.POST("/update-profile", s.UpdateProfile)
    
    acl.WithRouterOption(g, authz.WithRouterWhitelist("/register", "/login"))
    acl.WithRouterOption(g,
        authz.WithRouterPermission("/list", "read"),
        authz.WithRouterPermission("/delete", "delete"),
        authz.WithRouterPermission("/reset-password", "reset-password"),
    )
}

// Register
// @x-apifox-folder "用户"
// @Summary 用户注册
// @Tags    用户
// @Accept  json
// @Produce json
// @Param   body body dto.UserCreateParams true "请求参数"
// @Success 200 {object} types.ApiCreateResult
// @Router  /user/register [post]
func (s UserController) Register(ctx *gin.Context) {
    httpbinding.BindJSON[dto.UserCreateParams](ctx, func(params dto.UserCreateParams) (any, error) {
        params.IP = ctx.ClientIP()
        return s.UserSrv.Create(ctx, params)
    })
}

// Login
// @x-apifox-folder "用户"
// @Summary 用户登录
// @Tags    用户
// @Accept  json
// @Produce json
// @Param   body body dto.UserCreateParams true "请求参数"
// @Success 200 {object} types.ApiCreateResult
// @Router  /user/login [post]
func (s UserController) Login(ctx *gin.Context) {
    httpbinding.BindJSON[dto.UserCredentialParams](ctx, func(params dto.UserCredentialParams) (any, error) {
        u, err := s.UserSrv.FindByCredential(ctx, params)
        if err != nil {
            return nil, err
        }
        token, err := s.TokenSrv.Create(u.ID)
        if err != nil {
            return nil, err
        }
        return dto.UserAuthnResult{User: u, Token: token}, nil
    })
}

// Logout
// @x-apifox-folder "用户"
// @Summary 注销登录
// @Tags    用户
// @Accept  json
// @Produce json
// @Success 200 {object} types.ApiEmptyResult
// @Router  /user/logout [post]
func (s UserController) Logout(ctx *gin.Context) {
    tokenString := ctx.GetHeader("Authorization")
    if tokenString != "" {
        err := s.TokenSrv.Revoke(ctx, tokenString)
        if err != nil {
            s.Infra.Log.App.Warn("Token注销失败", zap.Error(err))
        }
    }
    erroz.OK.Write(ctx)
}

// ChangePassword
// @x-apifox-folder "用户"
// @Summary 修改密码
// @Tags    用户
// @Accept  json
// @Produce json
// @Param	body body dto.UserPasswdUpdateParams true "请求参数"
// @Success 200 {object} types.ApiEmptyResult
// @Router  /user/change-password [post]
func (s UserController) ChangePassword(ctx *gin.Context) {
    httpbinding.BindJSON[dto.UserPasswdUpdateParams](ctx, func(params dto.UserPasswdUpdateParams) (any, error) {
        u, err := contextutil.GetAuthorizedUser(ctx)
        if err != nil {
            return nil, err
        }
        return nil, s.UserSrv.ChangePassword(ctx, u, params)
    })
}

// ResetPassword
// @x-apifox-folder "用户"
// @Summary 重置密码
// @Tags    用户
// @Accept  json
// @Produce json
// @Param	body body dto.UserPasswdResetParams true "请求参数"
// @Success 200 {object} types.ApiEmptyResult
// @Router  /user/reset-password [post]
func (s UserController) ResetPassword(ctx *gin.Context) {
    httpbinding.BindJSON[dto.UserPasswdResetParams](ctx, func(params dto.UserPasswdResetParams) (any, error) {
        userInfo, err := s.UserSrv.FindByID(ctx, params.ID)
        if err != nil {
            return nil, err
        }
        return nil, s.UserSrv.ResetPassword(ctx, userInfo.ID, params.Password)
    })
}

// List
// @x-apifox-folder "用户"
// @Summary 查看用户列表
// @Tags    用户
// @Accept  json
// @Produce json
// @Param	body body dto.UserListParams true "请求参数"
// @Success 200 {object} swagger.UserList
// @Router  /user/list [post]
func (s UserController) List(ctx *gin.Context) {
    httpbinding.BindJSON[dto.UserListParams](ctx, func(params dto.UserListParams) (any, error) {
        return s.UserSrv.List(ctx, params)
    })
}

// Detail
// @x-apifox-folder "用户"
// @Summary 查看用户资料
// @Tags    用户
// @Accept  json
// @Produce json
// @Param	body body types.ResourceID true "请求参数"
// @Success 200 {object} swagger.User
// @Router  /user/detail [post]
func (s UserController) Detail(ctx *gin.Context) {
    httpbinding.BindJSON[types.ResourceID](ctx, func(params types.ResourceID) (any, error) {
        return s.UserSrv.FindByID(ctx, params.ID)
    })
}

// Delete
// @x-apifox-folder "用户"
// @Summary 删除用户
// @Tags    用户
// @Accept  json
// @Produce json
// @Param	body body types.ResourceID true "请求参数"
// @Success 200 {object} types.ApiEmptyResult
// @Router  /user/delete [post]
func (s UserController) Delete(ctx *gin.Context) {
    httpbinding.BindJSON[types.ResourceID](ctx, func(params types.ResourceID) (any, error) {
        return nil, s.UserSrv.Delete(ctx, params.ID)
    })
}

// UpdateProfile
// @x-apifox-folder "用户"
// @Summary 更新用户资料
// @Tags    用户
// @Accept  json
// @Produce json
// @Param	body body dto.UserProfile true "请求参数"
// @Success 200 {object} types.ApiEmptyResult
// @Router  /user/update-profile [post]
func (s UserController) UpdateProfile(ctx *gin.Context) {
    httpbinding.BindJSON[dto.UserProfile](ctx, func(params dto.UserProfile) (any, error) {
        u, err := contextutil.GetAuthorizedUser(ctx)
        if err != nil {
            return nil, err
        }
        params.UserID = u.ID
        return nil, s.UserSrv.UpdateProfile(ctx, params)
    })
}
