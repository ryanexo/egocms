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

func (c RoleController) Create(ctx *gin.Context) {
    httpbinding.BindJSON[dto.RoleCreateParams](ctx, func(params dto.RoleCreateParams) (any, error) {
        return c.srv.Role.Create(ctx, params)
    })
}

func (c RoleController) Update(ctx *gin.Context) {
    httpbinding.BindJSON[dto.RoleUpdateParams](ctx, func(params dto.RoleUpdateParams) (any, error) {
        return nil, c.srv.Role.Update(ctx, params)
    })
}

func (c RoleController) Delete(ctx *gin.Context) {
    httpbinding.BindJSON[dto.ResourceID](ctx, func(params dto.ResourceID) (any, error) {
        return nil, c.srv.Role.Delete(ctx, params.ID)
    })
}

func (c RoleController) List(ctx *gin.Context) {
    httpbinding.BindJSON[dto.RoleListParams](ctx, func(params dto.RoleListParams) (any, error) {
        return c.srv.Role.List(ctx, params)
    })
}

func (c RoleController) Detail(ctx *gin.Context) {
    httpbinding.BindJSON[dto.ResourceID](ctx, func(params dto.ResourceID) (any, error) {
        return c.srv.Role.FindRoleByID(ctx, params.ID)
    })
}
