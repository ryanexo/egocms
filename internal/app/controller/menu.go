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

// List 菜单列表
// @x-apifox-folder "菜单"
// @Summary 菜单列表
// @Tags    菜单
// @Accept  json
// @Produce json
// @Param   body body dto.MenuListQueryParams true "请求参数"
// @Success 200 {object} erroz.Result{data=nil}
// @Router  /menu/list [post]
func (c MenuController) List(ctx *gin.Context) {
    httpbinding.BindJSON[dto.MenuListQueryParams](ctx, func(params dto.MenuListQueryParams) (any, error) {
        return c.srv.Menu.List(ctx, params)
    })
}

// Detail 查看菜单
// @x-apifox-folder "菜单"
// @Summary 查看菜单
// @Tags    菜单
// @Accept  json
// @Produce json
// @Param   body body dto.ResourceID true "请求参数"
// @Success 200 {object} erroz.Result{data=nil}
// @Router  /menu/detail [post]
func (c MenuController) Detail(ctx *gin.Context) {
    httpbinding.BindJSON[dto.ResourceID](ctx, func(params dto.ResourceID) (any, error) {
        return c.srv.Menu.FindByID(ctx, params.ID)
    })
}

// Create 创建菜单
// @x-apifox-folder "菜单"
// @Summary 创建菜单
// @Tags    菜单
// @Accept  json
// @Produce json
// @Param   body body dto.MenuCreateParams true "请求参数"
// @Success 200 {object} erroz.Result{data=nil}
// @Router  /menu/create [post]
func (c MenuController) Create(ctx *gin.Context) {
    httpbinding.BindJSON[dto.MenuCreateParams](ctx, func(params dto.MenuCreateParams) (any, error) {
        return c.srv.Menu.Create(ctx, params)
    })
}

// Move 移动菜单
// @x-apifox-folder "菜单"
// @Summary 移动菜单
// @Tags    菜单
// @Accept  json
// @Produce json
// @Param   body body dto.MenuMoveParams true "请求参数"
// @Success 200 {object} erroz.Result{data=nil}
// @Router  /menu/move [post]
func (c MenuController) Move(ctx *gin.Context) {
    httpbinding.BindJSON[dto.MenuMoveParams](ctx, func(params dto.MenuMoveParams) (any, error) {
        return nil, c.srv.Menu.Move(ctx, params.ID, params.TargetID)
    })
}

// Update 更新菜单
// @x-apifox-folder "菜单"
// @Summary 更新菜单
// @Tags    菜单
// @Accept  json
// @Produce json
// @Param   body body dto.MenuUpdateParams true "请求参数"
// @Success 200 {object} erroz.Result{data=nil}
// @Router  /menu/update [post]
func (c MenuController) Update(ctx *gin.Context) {
    httpbinding.BindJSON[dto.MenuUpdateParams](ctx, func(params dto.MenuUpdateParams) (any, error) {
        err := c.srv.Menu.Update(ctx, params)
        return nil, err
    })
}

// Delete 删除菜单
// @x-apifox-folder "菜单"
// @Summary 删除菜单
// @Tags    菜单
// @Accept  json
// @Produce json
// @Param   body body dto.ResourceID true "请求参数"
// @Success 200 {object} erroz.Result{data=nil}
// @Router  /menu/delete [post]
func (c MenuController) Delete(ctx *gin.Context) {
    httpbinding.BindJSON[dto.ResourceID](ctx, func(params dto.ResourceID) (any, error) {
        err := c.srv.Menu.Delete(ctx, params.ID)
        return nil, err
    })
}
