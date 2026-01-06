package controller

import (
	"dpcms/internal/app/controller/internal/httpbinding"
	"dpcms/internal/app/dto"
	"dpcms/internal/app/erroz"
	"dpcms/internal/app/middleware/authz"
	"dpcms/internal/app/service"
	"dpcms/internal/app/util/contextutil"
	"dpcms/internal/infra"
	"dpcms/internal/infra/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserController struct {
	services *service.Services
	log      *logger.Logger
}

func NewUserController(s *service.Services, i *infra.Infra) UserController {
	return UserController{services: s, log: i.Log}
}

func (c UserController) setup(server *gin.Engine) {
	acl := authz.NewWithRBAC(c.services, "user")

	g := server.Group("/user", acl.Middleware())
	g.POST("/register", c.Register)
	g.POST("/login", c.Login)
	g.POST("/logout", c.Logout)
	g.POST("/update-password", c.ChangePassword)
	g.POST("/reset-password", c.ResetPassword)
	g.POST("/list", c.List)
	g.POST("/detail", c.Detail)
	g.POST("/delete", c.Delete)
	g.POST("/update-profile", c.UpdateProfile)

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
// @Success 200 {object} swaggertype.CreateResult
// @Router  /user/register [post]
func (c UserController) Register(ctx *gin.Context) {
	httpbinding.BindJSON[dto.UserCreateParams](ctx, func(params dto.UserCreateParams) (any, error) {
		params.IP = ctx.ClientIP()
		return c.services.User.Create(ctx, params)
	})
}

// Login
// @x-apifox-folder "用户"
// @Summary 用户登录
// @Tags    用户
// @Accept  json
// @Produce json
// @Param   body body dto.UserCreateParams true "请求参数"
// @Success 200 {object} swaggertype.CreateResult
// @Router  /user/login [post]
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

// Logout
// @x-apifox-folder "用户"
// @Summary 注销登录
// @Tags    用户
// @Accept  json
// @Produce json
// @Success 200 {object} swaggertype.EmptyResult
// @Router  /user/logout [post]
func (c UserController) Logout(ctx *gin.Context) {
	tokenString := ctx.GetHeader("Authorization")
	if tokenString != "" {
		err := c.services.Token.Revoke(ctx, tokenString)
		if err != nil {
			c.log.App.Warn("Token注销失败", zap.Error(err))
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
// @Success 200 {object} swaggertype.EmptyResult
// @Router  /user/change-password [post]
func (c UserController) ChangePassword(ctx *gin.Context) {
	httpbinding.BindJSON[dto.UserPasswdUpdateParams](ctx, func(params dto.UserPasswdUpdateParams) (any, error) {
		u, err := contextutil.GetAuthorizedUser(ctx)
		if err != nil {
			return nil, err
		}
		return nil, c.services.User.ChangePassword(ctx, u, params)
	})
}

// ResetPassword
// @x-apifox-folder "用户"
// @Summary 重置密码
// @Tags    用户
// @Accept  json
// @Produce json
// @Param	body body dto.UserPasswdResetParams true "请求参数"
// @Success 200 {object} swaggertype.EmptyResult
// @Router  /user/reset-password [post]
func (c UserController) ResetPassword(ctx *gin.Context) {
	httpbinding.BindJSON[dto.UserPasswdResetParams](ctx, func(params dto.UserPasswdResetParams) (any, error) {
		userInfo, err := c.services.User.FindByID(ctx, params.ID)
		if err != nil {
			return nil, err
		}
		return nil, c.services.User.ResetPassword(ctx, userInfo.ID, params.Password)
	})
}

// List
// @x-apifox-folder "用户"
// @Summary 查看用户列表
// @Tags    用户
// @Accept  json
// @Produce json
// @Param	body body dto.UserListParams true "请求参数"
// @Success 200 {object} swaggertype.UserList
// @Router  /user/list [post]
func (c UserController) List(ctx *gin.Context) {
	httpbinding.BindJSON[dto.UserListParams](ctx, func(params dto.UserListParams) (any, error) {
		return c.services.User.List(ctx, params)
	})
}

// Detail
// @x-apifox-folder "用户"
// @Summary 查看用户资料
// @Tags    用户
// @Accept  json
// @Produce json
// @Param	body body dto.ResourceID true "请求参数"
// @Success 200 {object} swaggertype.User
// @Router  /user/detail [post]
func (c UserController) Detail(ctx *gin.Context) {
	httpbinding.BindJSON[dto.ResourceID](ctx, func(params dto.ResourceID) (any, error) {
		return c.services.User.FindByID(ctx, params.ID)
	})
}

// Delete
// @x-apifox-folder "用户"
// @Summary 删除用户
// @Tags    用户
// @Accept  json
// @Produce json
// @Param	body body dto.ResourceID true "请求参数"
// @Success 200 {object} swaggertype.EmptyResult
// @Router  /user/delete [post]
func (c UserController) Delete(ctx *gin.Context) {
	httpbinding.BindJSON[dto.ResourceID](ctx, func(params dto.ResourceID) (any, error) {
		return nil, c.services.User.Delete(ctx, params.ID)
	})
}

// UpdateProfile
// @x-apifox-folder "用户"
// @Summary 更新用户资料
// @Tags    用户
// @Accept  json
// @Produce json
// @Param	body body dto.UserProfile true "请求参数"
// @Success 200 {object} swaggertype.EmptyResult
// @Router  /user/update-profile [post]
func (c UserController) UpdateProfile(ctx *gin.Context) {
	httpbinding.BindJSON[dto.UserProfile](ctx, func(params dto.UserProfile) (any, error) {
		u, err := contextutil.GetAuthorizedUser(ctx)
		if err != nil {
			return nil, err
		}
		params.UserID = u.ID
		return nil, c.services.User.UpdateProfile(ctx, params)
	})
}
