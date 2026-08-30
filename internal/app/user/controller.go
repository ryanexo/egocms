package user

import (
    "context"
    "strings"
    
    "cms/internal/app/user/api"
    "cms/internal/app/user/model"
    "cms/internal/httpx"
    "cms/internal/middleware/authz"
    "cms/internal/public/apitype"
    "cms/internal/public/jsontype"
    
    "github.com/gin-gonic/gin"
)

type TokenManager interface {
    Create(userID jsontype.SafeUint64) (string, error)
    Revoke(ctx context.Context, token string) error
}

type Controller struct {
    service *Service
    tokens  TokenManager
    auth    *authz.Factory
}

func NewUserController(service *Service, tokens TokenManager, auth *authz.Factory) *Controller {
    return &Controller{service: service, tokens: tokens, auth: auth}
}

var _ httpx.Route = (*Controller)(nil)

func (controller *Controller) Setup(router httpx.Router) {
    auth := controller.auth.Resource("user")
    group := router.Group("/user")
    group.POST("/register", auth.Public().Wrap(httpx.HandlerFunc(controller.Register)))
    group.POST("/login", auth.Public().Wrap(httpx.HandlerFunc(controller.Login)))
    group.POST("/logout", auth.Wrap(httpx.HandlerFunc(controller.Logout)))
    group.POST("/update-password", auth.Wrap(httpx.HandlerFunc(controller.ChangePassword)))
    group.POST("/reset-password", auth.Permission("reset-password").Wrap(httpx.HandlerFunc(controller.ResetPassword)))
    group.POST("/list", auth.Permission("read").Wrap(httpx.HandlerFunc(controller.List)))
    group.POST("/detail", auth.Wrap(httpx.HandlerFunc(controller.Detail)))
    group.POST("/delete", auth.Permission("delete").Wrap(httpx.HandlerFunc(controller.Delete)))
    group.POST("/update-profile", auth.Wrap(httpx.HandlerFunc(controller.UpdateProfile)))
}

// Register 用户注册
// @x-apifox-folder "用户"
// @Summary 用户注册
// @Tags 用户
// @Accept json
// @Produce json
// @Param body body api.UserCreateParams true "请求参数"
// @Success 200 {object} apitype.ApiCreateResult
// @Router /user/register [post]
func (controller *Controller) Register(ctx *gin.Context) error {
    return httpx.BindJSON[api.UserCreateParams](ctx, func(params api.UserCreateParams) (any, error) {
        params.IP = ctx.ClientIP()
        return controller.service.Create(ctx, params)
    })
}

// Login 用户登录
// @x-apifox-folder "用户"
// @Summary 用户登录
// @Tags 用户
// @Accept json
// @Produce json
// @Param body body api.UserCredentialParams true "请求参数"
// @Success 200 {object} api.ApiUserLogin
// @Router /user/login [post]
func (controller *Controller) Login(ctx *gin.Context) error {
    return httpx.BindJSON[api.UserCredentialParams](ctx, func(params api.UserCredentialParams) (any, error) {
        user, err := controller.service.FindByCredential(ctx, params)
        if err != nil {
            return nil, err
        }
        token, err := controller.tokens.Create(user.ID)
        if err != nil {
            return nil, err
        }
        return api.UserAuthnResult{User: user, Token: token}, nil
    })
}

// Logout 注销登录
// @x-apifox-folder "用户"
// @Security ApiKeyAuth
// @Summary 注销登录
// @Tags 用户
// @Produce json
// @Success 200 {object} apitype.ApiEmptyResult
// @Router /user/logout [post]
func (controller *Controller) Logout(ctx *gin.Context) error {
    credential := ctx.GetHeader("Authorization")
    token, found := strings.CutPrefix(credential, "Bearer ")
    if !found || token == "" {
        return authz.ErrAuthorized
    }
    if err := controller.tokens.Revoke(ctx, token); err != nil {
        return err
    }
    httpx.JSON(ctx, httpx.OK)
    return nil
}

// ChangePassword 修改密码
// @x-apifox-folder "用户"
// @Security ApiKeyAuth
// @Summary 修改密码
// @Tags 用户
// @Accept json
// @Produce json
// @Param body body api.UserPasswordUpdateParams true "请求参数"
// @Success 200 {object} apitype.ApiEmptyResult
// @Router /user/update-password [post]
func (controller *Controller) ChangePassword(ctx *gin.Context) error {
    return httpx.BindJSON[api.UserPasswordUpdateParams](ctx, func(params api.UserPasswordUpdateParams) (any, error) {
        current, err := authz.CurrentUser[*model.User](ctx)
        if err != nil {
            return nil, err
        }
        return nil, controller.service.ChangePassword(ctx, current, params)
    })
}

// ResetPassword 重置密码
// @x-apifox-folder "用户"
// @Security ApiKeyAuth
// @Summary 重置密码
// @Tags 用户
// @Accept json
// @Produce json
// @Param body body api.UserPasswordResetParams true "请求参数"
// @Success 200 {object} apitype.ApiEmptyResult
// @Router /user/reset-password [post]
func (controller *Controller) ResetPassword(ctx *gin.Context) error {
    return httpx.BindJSON[api.UserPasswordResetParams](ctx, func(params api.UserPasswordResetParams) (any, error) {
        return nil, controller.service.ResetPassword(ctx, params.ID.Uint64(), params.Password)
    })
}

// List 用户列表
// @x-apifox-folder "用户"
// @Security ApiKeyAuth
// @Summary 用户列表
// @Tags 用户
// @Accept json
// @Produce json
// @Param body body api.UserListParams true "请求参数"
// @Success 200 {object} api.ApiUserList
// @Router /user/list [post]
func (controller *Controller) List(ctx *gin.Context) error {
    return httpx.BindJSON[api.UserListParams](ctx, func(params api.UserListParams) (any, error) {
        return controller.service.List(ctx, params)
    })
}

// Detail 查看用户资料
// @x-apifox-folder "用户"
// @Security ApiKeyAuth
// @Summary 查看用户资料
// @Tags 用户
// @Accept json
// @Produce json
// @Param body body apitype.ResourceID true "请求参数"
// @Success 200 {object} api.ApiUser
// @Router /user/detail [post]
func (controller *Controller) Detail(ctx *gin.Context) error {
    return httpx.BindJSON[apitype.ResourceID](ctx, func(params apitype.ResourceID) (any, error) {
        return controller.service.FindByID(ctx, params.ID.Uint64())
    })
}

// Delete 删除用户
// @x-apifox-folder "用户"
// @Security ApiKeyAuth
// @Summary 删除用户
// @Tags 用户
// @Accept json
// @Produce json
// @Param body body apitype.ResourceID true "请求参数"
// @Success 200 {object} apitype.ApiEmptyResult
// @Router /user/delete [post]
func (controller *Controller) Delete(ctx *gin.Context) error {
    return httpx.BindJSON[apitype.ResourceID](ctx, func(params apitype.ResourceID) (any, error) {
        return nil, controller.service.Delete(ctx, params.ID.Uint64())
    })
}

// UpdateProfile 更新用户资料
// @x-apifox-folder "用户"
// @Security ApiKeyAuth
// @Summary 更新用户资料
// @Tags 用户
// @Accept json
// @Produce json
// @Param body body api.UserProfileParams true "请求参数"
// @Success 200 {object} apitype.ApiEmptyResult
// @Router /user/update-profile [post]
func (controller *Controller) UpdateProfile(ctx *gin.Context) error {
    return httpx.BindJSON[api.UserProfileParams](ctx, func(params api.UserProfileParams) (any, error) {
        current, err := authz.CurrentUser[*model.User](ctx)
        if err != nil {
            return nil, err
        }
        params.UserID = jsontype.SafeUint64(current.ID)
        return nil, controller.service.UpdateProfile(ctx, params)
    })
}
