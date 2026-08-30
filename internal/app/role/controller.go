package role

import (
    "cms/internal/app/role/api"
    "cms/internal/httpx"
    "cms/internal/middleware/authz"
    "cms/internal/public/apitype"
    
    "github.com/gin-gonic/gin"
)

type Controller struct {
    service *Service
    auth    *authz.Factory
}

func NewRoleController(service *Service, auth *authz.Factory) *Controller {
    return &Controller{service: service, auth: auth}
}

func (controller *Controller) Setup(router httpx.Router) {
    auth := controller.auth.Resource("role")
    group := router.Group("/role")
    group.POST("/list", auth.Permission("read").Wrap(httpx.HandlerFunc(controller.List)))
    group.POST("/create", auth.Permission("create").Wrap(httpx.HandlerFunc(controller.Create)))
    group.POST("/update", auth.Permission("update").Wrap(httpx.HandlerFunc(controller.Update)))
    group.POST("/delete", auth.Permission("delete").Wrap(httpx.HandlerFunc(controller.Delete)))
    group.POST("/detail", auth.Permission("read").Wrap(httpx.HandlerFunc(controller.Detail)))
}

// List 角色列表
// @x-apifox-folder "角色"
// @Security ApiKeyAuth
// @Summary 角色列表
// @Tags 角色
// @Accept json
// @Produce json
// @Param body body api.RoleListParams true "请求参数"
// @Success 200 {object} api.ApiRoleList
// @Router /role/list [post]
func (controller *Controller) List(ctx *gin.Context) error {
    return httpx.BindJSON[api.RoleListParams](ctx, func(params api.RoleListParams) (any, error) {
        return controller.service.List(ctx, params)
    })
}

// Detail 查看角色
// @x-apifox-folder "角色"
// @Security ApiKeyAuth
// @Summary 查看角色
// @Tags 角色
// @Accept json
// @Produce json
// @Param body body apitype.ResourceID true "请求参数"
// @Success 200 {object} api.ApiRole
// @Router /role/detail [post]
func (controller *Controller) Detail(ctx *gin.Context) error {
    return httpx.BindJSON[apitype.ResourceID](ctx, func(params apitype.ResourceID) (any, error) {
        return controller.service.FindByID(ctx, params.ID.Uint64())
    })
}

// Create 创建角色
// @x-apifox-folder "角色"
// @Security ApiKeyAuth
// @Summary 创建角色
// @Tags 角色
// @Accept json
// @Produce json
// @Param body body api.RoleCreateParams true "请求参数"
// @Success 200 {object} apitype.ApiCreateResult
// @Router /role/create [post]
func (controller *Controller) Create(ctx *gin.Context) error {
    return httpx.BindJSON[api.RoleCreateParams](ctx, func(params api.RoleCreateParams) (any, error) {
        return controller.service.Create(ctx, params)
    })
}

// Update 更新角色
// @x-apifox-folder "角色"
// @Security ApiKeyAuth
// @Summary 更新角色
// @Tags 角色
// @Accept json
// @Produce json
// @Param body body api.RoleUpdateParams true "请求参数"
// @Success 200 {object} apitype.ApiEmptyResult
// @Router /role/update [post]
func (controller *Controller) Update(ctx *gin.Context) error {
    return httpx.BindJSON[api.RoleUpdateParams](ctx, func(params api.RoleUpdateParams) (any, error) {
        return nil, controller.service.Update(ctx, params)
    })
}

// Delete 删除角色
// @x-apifox-folder "角色"
// @Security ApiKeyAuth
// @Summary 删除角色
// @Tags 角色
// @Accept json
// @Produce json
// @Param body body apitype.ResourceID true "请求参数"
// @Success 200 {object} apitype.ApiEmptyResult
// @Router /role/delete [post]
func (controller *Controller) Delete(ctx *gin.Context) error {
    return httpx.BindJSON[apitype.ResourceID](ctx, func(params apitype.ResourceID) (any, error) {
        return nil, controller.service.Delete(ctx, params.ID.Uint64())
    })
}
