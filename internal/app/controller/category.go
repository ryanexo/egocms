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

func (c CategoryController) List(ctx *gin.Context) {
    httpbinding.BindJSON[dto.CategoryListParams](ctx, func(params dto.CategoryListParams) (any, error) {
        return c.srv.Category.List(ctx, params)
    })
}

func (c CategoryController) Detail(ctx *gin.Context) {
    httpbinding.BindJSON[dto.QueryByResourceID](ctx, func(params dto.QueryByResourceID) (any, error) {
        return c.srv.Category.FindByID(ctx, params.ID)
    })
}

func (c CategoryController) Create(ctx *gin.Context) {
    httpbinding.BindJSON[dto.CategoryCreateParams](ctx, func(params dto.CategoryCreateParams) (any, error) {
        return c.srv.Category.Create(ctx, params)
    })
}

func (c CategoryController) Move(ctx *gin.Context) {
    httpbinding.BindJSON[dto.CategoryMoveParams](ctx, func(params dto.CategoryMoveParams) (any, error) {
        return nil, c.srv.Category.Move(ctx, params.ID, params.TargetID)
    })
}

func (c CategoryController) Update(ctx *gin.Context) {
    httpbinding.BindJSON[dto.CategoryUpdateParams](ctx, func(params dto.CategoryUpdateParams) (any, error) {
        err := c.srv.Category.Update(ctx, params)
        return nil, err
    })
}

func (c CategoryController) Delete(ctx *gin.Context) {
    httpbinding.BindJSON[dto.QueryByResourceID](ctx, func(params dto.QueryByResourceID) (any, error) {
        err := c.srv.Category.Delete(ctx, params.ID)
        return nil, err
    })
}
