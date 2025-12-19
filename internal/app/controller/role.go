package controller

import (
    `dpcms/internal/app/controller/internal/common`
    `dpcms/internal/app/middleware/authz`
    `dpcms/internal/app/service`
    `dpcms/internal/app/service/srvparams`
    `dpcms/internal/infra`
    
    `github.com/gin-gonic/gin`
)

type RoleController struct {
    srv   *service.Services
    infra *infra.Infra
}

func NewRoleController(srv *service.Services, infra *infra.Infra) RoleController {
    return RoleController{srv: srv, infra: infra}
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
    common.BindJSON[srvparams.RoleCreateParams](ctx, func(params srvparams.RoleCreateParams) (any, error) {
        return c.srv.Role.Create(ctx, params)
    })
}

func (c RoleController) Update(ctx *gin.Context) {
    common.BindJSON[srvparams.RoleUpdateParams](ctx, func(params srvparams.RoleUpdateParams) (any, error) {
        return nil, c.srv.Role.Update(ctx, params)
    })
}

func (c RoleController) Delete(ctx *gin.Context) {
    common.BindJSON[srvparams.QueryByResourceID](ctx, func(params srvparams.QueryByResourceID) (any, error) {
        return nil, c.srv.Role.Delete(ctx, params.ID)
    })
}

func (c RoleController) List(ctx *gin.Context) {
    common.BindJSON[srvparams.RoleListParams](ctx, func(params srvparams.RoleListParams) (any, error) {
        return c.srv.Role.List(ctx, params)
    })
}

func (c RoleController) Detail(ctx *gin.Context) {
    common.BindJSON[srvparams.QueryByResourceID](ctx, func(params srvparams.QueryByResourceID) (any, error) {
        return c.srv.Role.FindRoleById(ctx, params.ID)
    })
}
