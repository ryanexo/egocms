package controller

import (
    `dpcms/internal/app/controller/internal/httpbinding`
    `dpcms/internal/app/dto`
    `dpcms/internal/app/middleware/authz`
    `dpcms/internal/app/service`
    
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
    httpbinding.BindJSON[dto.MenuListQueryParams](ctx, func(params dto.MenuListQueryParams) (any, error) {
        return c.srv.Menu.List(ctx, params)
    })
}

func (c MenuController) Detail(ctx *gin.Context) {
    httpbinding.BindJSON[dto.ResourceID](ctx, func(params dto.ResourceID) (any, error) {
        return c.srv.Menu.FindByID(ctx, params.ID)
    })
}

func (c MenuController) Create(ctx *gin.Context) {
    httpbinding.BindJSON[dto.MenuCreateParams](ctx, func(params dto.MenuCreateParams) (any, error) {
        return c.srv.Menu.Create(ctx, params)
    })
}

func (c MenuController) Move(ctx *gin.Context) {
    httpbinding.BindJSON[dto.MenuMoveParams](ctx, func(params dto.MenuMoveParams) (any, error) {
        return nil, c.srv.Menu.Move(ctx, params.ID, params.TargetID)
    })
}

func (c MenuController) Update(ctx *gin.Context) {
    httpbinding.BindJSON[dto.MenuUpdateParams](ctx, func(params dto.MenuUpdateParams) (any, error) {
        err := c.srv.Menu.Update(ctx, params)
        return nil, err
    })
}

func (c MenuController) Delete(ctx *gin.Context) {
    httpbinding.BindJSON[dto.ResourceID](ctx, func(params dto.ResourceID) (any, error) {
        err := c.srv.Menu.Delete(ctx, params.ID)
        return nil, err
    })
}
