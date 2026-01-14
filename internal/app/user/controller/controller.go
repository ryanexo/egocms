package controller

import (
    token `cms/internal/app/token/service`
    `cms/internal/app/user/internal/dto`
    user `cms/internal/app/user/service`
    "cms/internal/erroz"
    `cms/internal/httpserver`
    `cms/internal/infra/logger`
    `cms/internal/middleware/authz`
    `cms/internal/types`
    `cms/internal/util/contextutil`
    `cms/internal/util/httpbinding`
    
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type UserController struct {
    userSrv  *user.UserService
    tokenSrv *token.TokenService
    logger   *logger.Logger
    auth     *authz.Factory
}

func NewUserController(
    userSrv *user.UserService,
    tokenSrv *token.TokenService,
    logger *logger.Logger,
    auth *authz.Factory,
) *UserController {
    return &UserController{
        userSrv:  userSrv,
        tokenSrv: tokenSrv,
        logger:   logger,
        auth:     auth,
    }
}

func (s UserController) Setup(router httpserver.Router) {
    acl := s.auth.AccessControl("user")
    
    g := router.Group("/user", acl.Middleware())
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
        return s.userSrv.Create(ctx, params)
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
        u, err := s.userSrv.FindByCredential(ctx, params)
        if err != nil {
            return nil, err
        }
        tokenStr, err := s.tokenSrv.Create(u.ID)
        if err != nil {
            return nil, err
        }
        return dto.UserAuthnResult{User: u, Token: tokenStr}, nil
    })
}

// Logout
// @x-apifox-folder "用户"
// @Security ApiKeyAuth
// @Summary 注销登录
// @Tags    用户
// @Accept  json
// @Produce json
// @Success 200 {object} types.ApiEmptyResult
// @Router  /user/logout [post]
func (s UserController) Logout(ctx *gin.Context) {
    tokenString := ctx.GetHeader("Authorization")
    if tokenString != "" {
        err := s.tokenSrv.Revoke(ctx, tokenString)
        if err != nil {
            s.logger.App.Warn("Token注销失败", zap.Error(err))
        }
    }
    erroz.OK.Write(ctx)
}

// ChangePassword
// @x-apifox-folder "用户"
// @Security ApiKeyAuth
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
        return nil, s.userSrv.ChangePassword(ctx, u, params)
    })
}

// ResetPassword
// @x-apifox-folder "用户"
// @Security ApiKeyAuth
// @Summary 重置密码
// @Tags    用户
// @Accept  json
// @Produce json
// @Param	body body dto.UserPasswdResetParams true "请求参数"
// @Success 200 {object} types.ApiEmptyResult
// @Router  /user/reset-password [post]
func (s UserController) ResetPassword(ctx *gin.Context) {
    httpbinding.BindJSON[dto.UserPasswdResetParams](ctx, func(params dto.UserPasswdResetParams) (any, error) {
        userInfo, err := s.userSrv.FindByID(ctx, params.ID)
        if err != nil {
            return nil, err
        }
        return nil, s.userSrv.ResetPassword(ctx, userInfo.ID, params.Password)
    })
}

// List
// @x-apifox-folder "用户"
// @Security ApiKeyAuth
// @Summary 查看用户列表
// @Tags    用户
// @Accept  json
// @Produce json
// @Param	body body dto.UserListParams true "请求参数"
// @Success 200 {object} dto.ApiUserList
// @Router  /user/list [post]
func (s UserController) List(ctx *gin.Context) {
    httpbinding.BindJSON[dto.UserListParams](ctx, func(params dto.UserListParams) (any, error) {
        return s.userSrv.List(ctx, params)
    })
}

// Detail
// @x-apifox-folder "用户"
// @Security ApiKeyAuth
// @Summary 查看用户资料
// @Tags    用户
// @Accept  json
// @Produce json
// @Param	body body types.ResourceID true "请求参数"
// @Success 200 {object} dto.ApiUser
// @Router  /user/detail [post]
func (s UserController) Detail(ctx *gin.Context) {
    httpbinding.BindJSON[types.ResourceID](ctx, func(params types.ResourceID) (any, error) {
        return s.userSrv.FindByID(ctx, params.ID)
    })
}

// Delete
// @x-apifox-folder "用户"
// @Security ApiKeyAuth
// @Summary 删除用户
// @Tags    用户
// @Accept  json
// @Produce json
// @Param	body body types.ResourceID true "请求参数"
// @Success 200 {object} types.ApiEmptyResult
// @Router  /user/delete [post]
func (s UserController) Delete(ctx *gin.Context) {
    httpbinding.BindJSON[types.ResourceID](ctx, func(params types.ResourceID) (any, error) {
        return nil, s.userSrv.Delete(ctx, params.ID)
    })
}

// UpdateProfile
// @x-apifox-folder "用户"
// @Security ApiKeyAuth
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
        return nil, s.userSrv.UpdateProfile(ctx, params)
    })
}
