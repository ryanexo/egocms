package controller

import (
    `cms/internal/app/menu/internal/dto`
    `cms/internal/app/menu/service`
    `cms/internal/httpserver`
    `cms/internal/middleware/authz`
    `cms/internal/types`
    `cms/internal/util/httpbinding`
    
    `github.com/gin-gonic/gin`
)

type MenuController struct {
    menuSrv *service.MenuService
    auth    *authz.Factory
}

func NewMenuController(
    menuSrv *service.MenuService,
    auth *authz.Factory,
) *MenuController {
    return &MenuController{menuSrv: menuSrv, auth: auth}
}

func (s MenuController) Setup(router httpserver.Router) {
    acl := s.auth.AccessControl("menu")
    
    g := router.Group("/menu", acl.Middleware())
    g.POST("/list", s.List)
    g.POST("/create", s.Create)
    g.POST("/update", s.Update)
    g.POST("/delete", s.Delete)
    g.POST("/detail", s.Detail)
    g.POST("/move", s.Move)
    
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
// @Security ApiKeyAuth
// @Summary 菜单列表
// @Tags    菜单
// @Accept  json
// @Produce json
// @Param   body body dto.MenuListQueryParams true "请求参数"
// @Success 200 {object} dto.ApiMenuList
// @Router  /menu/list [post]
func (s MenuController) List(ctx *gin.Context) {
    httpbinding.BindJSON[dto.MenuListQueryParams](ctx, func(params dto.MenuListQueryParams) (any, error) {
        return s.menuSrv.List(ctx, params)
    })
}

// Detail 查看菜单
// @x-apifox-folder "菜单"
// @Security ApiKeyAuth
// @Summary 查看菜单
// @Tags    菜单
// @Accept  json
// @Produce json
// @Param   body body types.ResourceID true "请求参数"
// @Success 200 {object} dto.ApiMenu
// @Router  /menu/detail [post]
func (s MenuController) Detail(ctx *gin.Context) {
    httpbinding.BindJSON[types.ResourceID](ctx, func(params types.ResourceID) (any, error) {
        return s.menuSrv.FindByID(ctx, params.ID)
    })
}

// Create 创建菜单
// @x-apifox-folder "菜单"
// @Security ApiKeyAuth
// @Summary 创建菜单
// @Tags    菜单
// @Accept  json
// @Produce json
// @Param   body body dto.MenuCreateParams true "请求参数"
// @Success 200 {object} types.ApiEmptyResult
// @Router  /menu/create [post]
func (s MenuController) Create(ctx *gin.Context) {
    httpbinding.BindJSON[dto.MenuCreateParams](ctx, func(params dto.MenuCreateParams) (any, error) {
        return s.menuSrv.Create(ctx, params)
    })
}

// Move 移动菜单
// @x-apifox-folder "菜单"
// @Security ApiKeyAuth
// @Summary 移动菜单
// @Tags    菜单
// @Accept  json
// @Produce json
// @Param   body body dto.MenuMoveParams true "请求参数"
// @Success 200 {object} types.ApiEmptyResult
// @Router  /menu/move [post]
func (s MenuController) Move(ctx *gin.Context) {
    httpbinding.BindJSON[dto.MenuMoveParams](ctx, func(params dto.MenuMoveParams) (any, error) {
        return nil, s.menuSrv.Move(ctx, params.ID, params.TargetID)
    })
}

// Update 更新菜单
// @x-apifox-folder "菜单"
// @Security ApiKeyAuth
// @Summary 更新菜单
// @Tags    菜单
// @Accept  json
// @Produce json
// @Param   body body dto.MenuUpdateParams true "请求参数"
// @Success 200 {object} types.ApiEmptyResult
// @Router  /menu/update [post]
func (s MenuController) Update(ctx *gin.Context) {
    httpbinding.BindJSON[dto.MenuUpdateParams](ctx, func(params dto.MenuUpdateParams) (any, error) {
        err := s.menuSrv.Update(ctx, params)
        return nil, err
    })
}

// Delete 删除菜单
// @x-apifox-folder "菜单"
// @Security ApiKeyAuth
// @Summary 删除菜单
// @Tags    菜单
// @Accept  json
// @Produce json
// @Param   body body types.ResourceID true "请求参数"
// @Success 200 {object} types.ApiEmptyResult
// @Router  /menu/delete [post]
func (s MenuController) Delete(ctx *gin.Context) {
    httpbinding.BindJSON[types.ResourceID](ctx, func(params types.ResourceID) (any, error) {
        err := s.menuSrv.Delete(ctx, params.ID)
        return nil, err
    })
}
