package controller

import (
	"dpcms/internal/app/category/internal/dto"
	"dpcms/internal/app/category/service"
	"dpcms/internal/middleware/authz"
	"dpcms/internal/types"
	"dpcms/internal/util/httpbinding"

	_ "dpcms/internal/app/category/swagger"
	_ "dpcms/internal/util/httpbinding"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	CategorySrv *service.CategoryService
	Auth        *authz.Builder
}

func (s CategoryController) setup(engine *gin.Engine) {
	acl := s.Auth.AccessControl("category")

	g := engine.Group("/category", acl.Middleware())
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

// List 分类列表
// @x-apifox-folder "文章分类"
// @Summary 分类列表
// @Tags    分类
// @Accept  json
// @Produce json
// @Param   body body dto.CategoryListParams true "请求参数"
// @Success 200 {object} swagger.CategoryList
// @Router  /category/list [post]
func (s CategoryController) List(ctx *gin.Context) {
	httpbinding.BindJSON[dto.CategoryListParams](ctx, func(params dto.CategoryListParams) (any, error) {
		return s.CategorySrv.List(ctx, params)
	})
}

// Detail 查看分类
// @x-apifox-folder "文章分类"
// @Summary 查看分类
// @Tags    分类
// @Accept  json
// @Produce json
// @Param   body body types.ResourceID true "请求参数"
// @Success 200 {object} swagger.Category
// @Router  /category/detail [post]
func (s CategoryController) Detail(ctx *gin.Context) {
	httpbinding.BindJSON[types.ResourceID](ctx, func(params types.ResourceID) (any, error) {
		return s.CategorySrv.FindByID(ctx, params.ID)
	})
}

// Create 创建分类
// @x-apifox-folder "文章分类"
// @Summary 创建分类
// @Tags    分类
// @Accept  json
// @Produce json
// @Param   body body dto.CategoryCreateParams true "请求参数"
// @Success 200 {object} types.ApiCreateResult
// @Router  /category/create [post]
func (s CategoryController) Create(ctx *gin.Context) {
	httpbinding.BindJSON[dto.CategoryCreateParams](ctx, func(params dto.CategoryCreateParams) (any, error) {
		return s.CategorySrv.Create(ctx, params)
	})
}

// Move 移动分类
// @x-apifox-folder "文章分类"
// @Summary 移动分类
// @Tags    分类
// @Accept  json
// @Produce json
// @Param   body body dto.CategoryMoveParams true "请求参数"
// @Success 200 {object} types.ApiEmptyResult
// @Router  /category/move [post]
func (s CategoryController) Move(ctx *gin.Context) {
	httpbinding.BindJSON[dto.CategoryMoveParams](ctx, func(params dto.CategoryMoveParams) (any, error) {
		return nil, s.CategorySrv.Move(ctx, params.ID, params.TargetID)
	})
}

// Update 更新分类
// @x-apifox-folder "文章分类"
// @Summary 更新分类
// @Tags    分类
// @Accept  json
// @Produce json
// @Param   body body dto.CategoryUpdateParams true "请求参数"
// @Success 200 {object} types.ApiEmptyResult
// @Router  /category/update [post]
func (s CategoryController) Update(ctx *gin.Context) {
	httpbinding.BindJSON[dto.CategoryUpdateParams](ctx, func(params dto.CategoryUpdateParams) (any, error) {
		err := s.CategorySrv.Update(ctx, params)
		return nil, err
	})
}

// Delete 删除分类
// @x-apifox-folder "文章分类"
// @Summary 删除分类
// @Tags    分类
// @Accept  json
// @Produce json
// @Param   body body types.ResourceID true "请求参数"
// @Success 200 {object} types.ApiEmptyResult
// @Router  /category/delete [post]
func (s CategoryController) Delete(ctx *gin.Context) {
	httpbinding.BindJSON[types.ResourceID](ctx, func(params types.ResourceID) (any, error) {
		err := s.CategorySrv.Delete(ctx, params.ID)
		return nil, err
	})
}
