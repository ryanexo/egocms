package controller

import (
    `dpcms/internal/app/controller/internal/common`
    `dpcms/internal/app/middleware/authz`
    `dpcms/internal/app/service`
    `dpcms/internal/app/service/srvparams`
    
    `github.com/gin-gonic/gin`
)

type MenuController struct {
    srv *service.Services
}

func NewMenuController(srv *service.Services) MenuController {
    return MenuController{srv}
}

func (c MenuController) setup(engine *gin.Engine) {
    acl := authz.NewWithRBAC(c.srv, "menu")
    
    g := engine.Group("/menu", acl.Middleware())
    g.POST("/list", c.List)
    g.POST("/create", c.Create)
    g.POST("/update", c.Update)
    g.POST("/delete", c.Delete)
    g.POST("/detail", c.Detail)
    g.POST("/move", c.Move)
    
    acl.WithRouterOption(
        g,
        authz.WithRouterPermission("/list", "read"),
        authz.WithRouterPermission("/detail", "read"),
        authz.WithRouterPermission("/create", "create"),
        authz.WithRouterPermission("/update", "update"),
        authz.WithRouterPermission("/delete", "delete"),
    )
}

func (c MenuController) List(ctx *gin.Context) {
    common.BasicBind[srvparams.MenuListQueryParams](ctx, func(params srvparams.MenuListQueryParams) (any, error) {
        return c.srv.Menu.List(ctx, params)
    })
}

func (c MenuController) Detail(ctx *gin.Context) {
    common.BasicBind[srvparams.QueryByResourceID](ctx, func(params srvparams.QueryByResourceID) (any, error) {
        return c.srv.Menu.FindByID(ctx, params.ID)
    })
}

func (c MenuController) Create(ctx *gin.Context) {
    common.BasicBind[srvparams.MenuCreateParams](ctx, func(params srvparams.MenuCreateParams) (any, error) {
        return c.srv.Menu.Create(ctx, params)
    })
}

func (c MenuController) Move(ctx *gin.Context) {
    common.BasicBind[srvparams.MenuMoveParams](ctx, func(params srvparams.MenuMoveParams) (any, error) {
        return nil, c.srv.Menu.Move(ctx, params.ID, params.TargetID)
    })
}

func (c MenuController) Update(ctx *gin.Context) {
    common.BasicBind[srvparams.MenuUpdateParams](ctx, func(params srvparams.MenuUpdateParams) (any, error) {
        err := c.srv.Menu.Update(ctx, params)
        return nil, err
    })
}

func (c MenuController) Delete(ctx *gin.Context) {
    common.BasicBind[srvparams.QueryByResourceID](ctx, func(params srvparams.QueryByResourceID) (any, error) {
        err := c.srv.Menu.Delete(ctx, params.ID)
        return nil, err
    })
}
