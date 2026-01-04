package controller

import (
    `dpcms/internal/app/controller/internal/httpbinding`
    `dpcms/internal/app/dto`
    `dpcms/internal/app/middleware/authz`
    `dpcms/internal/app/service`
    
    `github.com/gin-gonic/gin`
)

type CategoryController struct {
    srv *service.Services
}

func NewCategoryController(srv *service.Services) CategoryController {
    return CategoryController{srv}
}

func (c CategoryController) setup(engine *gin.Engine) {
    acl := authz.NewWithRBAC(c.srv, "category")
    
    g := engine.Group("/category", acl.Middleware())
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

// List 分类列表
// @x-apifox-folder "内容分类"
// @Summary 分类列表
// @Tags    分类
// @Accept  json
// @Produce json
// @Param   body body dto.CategoryListParams true "请求参数"
// @Success 200 {object} swaggertype.CategoryList
// @Router  /category/list [post]
func (c CategoryController) List(ctx *gin.Context) {
    httpbinding.BindJSON[dto.CategoryListParams](ctx, func(params dto.CategoryListParams) (any, error) {
        return c.srv.Category.List(ctx, params)
    })
}

// Detail 查看分类
// @x-apifox-folder "内容分类"
// @Summary 查看分类
// @Tags    分类
// @Accept  json
// @Produce json
// @Param   body body dto.ResourceID true "请求参数"
// @Success 200 {object} swaggertype.Category
// @Router  /category/detail [post]
func (c CategoryController) Detail(ctx *gin.Context) {
    httpbinding.BindJSON[dto.ResourceID](ctx, func(params dto.ResourceID) (any, error) {
        return c.srv.Category.FindByID(ctx, params.ID)
    })
}

// Create 创建分类
// @x-apifox-folder "内容分类"
// @Summary 创建分类
// @Tags    分类
// @Accept  json
// @Produce json
// @Param   body body dto.CategoryCreateParams true "请求参数"
// @Success 200 {object} swaggertype.Category
// @Router  /category/create [post]
func (c CategoryController) Create(ctx *gin.Context) {
    httpbinding.BindJSON[dto.CategoryCreateParams](ctx, func(params dto.CategoryCreateParams) (any, error) {
        return c.srv.Category.Create(ctx, params)
    })
}

// Move 移动分类
// @x-apifox-folder "内容分类"
// @Summary 移动分类
// @Tags    分类
// @Accept  json
// @Produce json
// @Param   body body dto.CategoryMoveParams true "请求参数"
// @Success 200 {object} erroz.Result{data=nil}
// @Router  /category/move [post]
func (c CategoryController) Move(ctx *gin.Context) {
    httpbinding.BindJSON[dto.CategoryMoveParams](ctx, func(params dto.CategoryMoveParams) (any, error) {
        return nil, c.srv.Category.Move(ctx, params.ID, params.TargetID)
    })
}

// Update 更新分类
// @x-apifox-folder "内容分类"
// @Summary 更新分类
// @Tags    分类
// @Accept  json
// @Produce json
// @Param   body body dto.CategoryUpdateParams true "请求参数"
// @Success 200 {object} erroz.Result{data=nil}
// @Router  /category/update [post]
func (c CategoryController) Update(ctx *gin.Context) {
    httpbinding.BindJSON[dto.CategoryUpdateParams](ctx, func(params dto.CategoryUpdateParams) (any, error) {
        err := c.srv.Category.Update(ctx, params)
        return nil, err
    })
}

// Delete 删除分类
// @x-apifox-folder "内容分类"
// @Summary 删除分类
// @Tags    分类
// @Accept  json
// @Produce json
// @Param   body body dto.ResourceID true "请求参数"
// @Success 200 {object} erroz.Result{data=nil}
// @Router  /category/delete [post]
func (c CategoryController) Delete(ctx *gin.Context) {
    httpbinding.BindJSON[dto.ResourceID](ctx, func(params dto.ResourceID) (any, error) {
        err := c.srv.Category.Delete(ctx, params.ID)
        return nil, err
    })
}
