package controller

import (
    `dpcms/internal/app/controller/internal/httpbinding`
    `dpcms/internal/app/dto`
    `dpcms/internal/app/middleware/authz`
    `dpcms/internal/app/service`
    
    `github.com/gin-gonic/gin`
)

type RoleController struct {
    srv *service.Services
}

func NewRoleController(srv *service.Services) RoleController {
    return RoleController{srv: srv}
}

func (c RoleController) setup(server *gin.Engine) {
    acl := authz.NewWithRBAC(c.srv, "role")
    
    g := server.Group("/role", acl.Middleware())
    g.POST("/create", c.Create)
    g.POST("/update", c.Update)
    g.POST("/delete", c.Delete)
    g.POST("/list", c.List)
    g.POST("/detail", c.Detail)
    
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
// @Summary 创建角色
// @Tags 角色
// @Accept json
// @Produce json
// @Param body body dto.RoleCreateParams true "请求参数"
// @Success 200 {object} swaggertype.CreateResult
// @Router /role/create [post]
func (c RoleController) Create(ctx *gin.Context) {
    httpbinding.BindJSON[dto.RoleCreateParams](ctx, func(params dto.RoleCreateParams) (any, error) {
        return c.srv.Role.Create(ctx, params)
    })
}

// Update
// @x-apifox-folder "角色"
// @Summary 更新角色
// @Tags 角色
// @Accept json
// @Produce json
// @Param body body dto.RoleUpdateParams true "请求参数"
// @Success 200 {object} swaggertype.EmptyResult
// @Router /role/update [post]
func (c RoleController) Update(ctx *gin.Context) {
    httpbinding.BindJSON[dto.RoleUpdateParams](ctx, func(params dto.RoleUpdateParams) (any, error) {
        return nil, c.srv.Role.Update(ctx, params)
    })
}

// Delete
// @x-apifox-folder "角色"
// @Summary 删除角色
// @Tags 角色
// @Accept json
// @Produce json
// @Param body body dto.ResourceID true "请求参数"
// @Success 200 {object} swaggertype.EmptyResult
// @Router /role/delete [post]
func (c RoleController) Delete(ctx *gin.Context) {
    httpbinding.BindJSON[dto.ResourceID](ctx, func(params dto.ResourceID) (any, error) {
        return nil, c.srv.Role.Delete(ctx, params.ID)
    })
}

// List
// @x-apifox-folder "角色"
// @Summary 角色列表
// @Tags 角色
// @Accept json
// @Produce json
// @Param body body dto.RoleListParams true "请求参数"
// @Success 200 {object} swaggertype.RoleList
// @Router /role/list [post]
func (c RoleController) List(ctx *gin.Context) {
    httpbinding.BindJSON[dto.RoleListParams](ctx, func(params dto.RoleListParams) (any, error) {
        return c.srv.Role.List(ctx, params)
    })
}

// Detail
// @x-apifox-folder "角色"
// @Summary 查看角色信息
// @Tags 角色
// @Accept json
// @Produce json
// @Param body body dto.ResourceID true "请求参数"
// @Success 200 {object} swaggertype.Role
// @Router /role/detail [post]
func (c RoleController) Detail(ctx *gin.Context) {
    httpbinding.BindJSON[dto.ResourceID](ctx, func(params dto.ResourceID) (any, error) {
        return c.srv.Role.FindByID(ctx, params.ID)
    })
}
