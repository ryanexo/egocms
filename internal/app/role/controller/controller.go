package controller

import (
    `cms/internal/app/role/internal/dto`
    `cms/internal/app/role/service`
    `cms/internal/httpserver`
    `cms/internal/middleware/authz`
    `cms/internal/types`
    `cms/internal/util/httpbinding`
    
    `github.com/gin-gonic/gin`
)

type RoleController struct {
    roleSrv *service.RoleService
    auth    *authz.Factory
}

func NewRoleController(
    roleSrv *service.RoleService,
    auth *authz.Factory,
) *RoleController {
    return &RoleController{roleSrv: roleSrv, auth: auth}
}

func (s RoleController) Setup(router httpserver.Router) {
    acl := s.auth.AccessControl("role")
    
    g := router.Group("/role", acl.Middleware())
    g.POST("/create", s.Create)
    g.POST("/update", s.Update)
    g.POST("/delete", s.Delete)
    g.POST("/list", s.List)
    g.POST("/detail", s.Detail)
    
    acl.WithRouterOption(
        g,
        authz.WithRouterPermission("/create", "create"),
        authz.WithRouterPermission("/update", "update"),
        authz.WithRouterPermission("/delete", "delete"),
        authz.WithRouterPermission("/list", "read"),
        authz.WithRouterPermission("/detail", "read"),
    )
}

// Create
// @x-apifox-folder "角色"
// @Security ApiKeyAuth
// @Summary 创建角色
// @Tags 角色
// @Accept json
// @Produce json
// @Param body body dto.RoleCreateParams true "请求参数"
// @Success 200 {object} types.ApiCreateResult
// @Router /role/create [post]
func (s RoleController) Create(ctx *gin.Context) {
    httpbinding.BindJSON[dto.RoleCreateParams](ctx, func(params dto.RoleCreateParams) (any, error) {
        return s.roleSrv.Create(ctx, params)
    })
}

// Update
// @x-apifox-folder "角色"
// @Security ApiKeyAuth
// @Summary 更新角色
// @Tags 角色
// @Accept json
// @Produce json
// @Param body body dto.RoleUpdateParams true "请求参数"
// @Success 200 {object} types.ApiEmptyResult
// @Router /role/update [post]
func (s RoleController) Update(ctx *gin.Context) {
    httpbinding.BindJSON[dto.RoleUpdateParams](ctx, func(params dto.RoleUpdateParams) (any, error) {
        return nil, s.roleSrv.Update(ctx, params)
    })
}

// Delete
// @x-apifox-folder "角色"
// @Security ApiKeyAuth
// @Summary 删除角色
// @Tags 角色
// @Accept json
// @Produce json
// @Param body body types.ResourceID true "请求参数"
// @Success 200 {object} types.ApiEmptyResult
// @Router /role/delete [post]
func (s RoleController) Delete(ctx *gin.Context) {
    httpbinding.BindJSON[types.ResourceID](ctx, func(params types.ResourceID) (any, error) {
        return nil, s.roleSrv.Delete(ctx, params.ID)
    })
}

// List
// @x-apifox-folder "角色"
// @Security ApiKeyAuth
// @Summary 角色列表
// @Tags 角色
// @Accept json
// @Produce json
// @Param body body dto.RoleListParams true "请求参数"
// @Success 200 {object} dto.ApiRoleList
// @Router /role/list [post]
func (s RoleController) List(ctx *gin.Context) {
    httpbinding.BindJSON[dto.RoleListParams](ctx, func(params dto.RoleListParams) (any, error) {
        return s.roleSrv.List(ctx, params)
    })
}

// Detail
// @x-apifox-folder "角色"
// @Security ApiKeyAuth
// @Summary 查看角色信息
// @Tags 角色
// @Accept json
// @Produce json
// @Param body body types.ResourceID true "请求参数"
// @Success 200 {object} dto.ApiRole
// @Router /role/detail [post]
func (s RoleController) Detail(ctx *gin.Context) {
    httpbinding.BindJSON[types.ResourceID](ctx, func(params types.ResourceID) (any, error) {
        return s.roleSrv.FindByID(ctx, params.ID)
    })
}
