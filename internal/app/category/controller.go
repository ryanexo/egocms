package category

import (
    "cms/internal/app/category/api"
    "cms/internal/httpx"
    "cms/internal/middleware/authz"
    "cms/internal/public/apitype"
    
    "github.com/gin-gonic/gin"
)

type CategoryController struct {
    service *Service
    auth    *authz.Factory
}

func NewCategoryController(service *Service, auth *authz.Factory) *CategoryController {
    return &CategoryController{service: service, auth: auth}
}

func (s *CategoryController) Setup(router httpx.Router) {
    auth := s.auth.Resource("category")
    group := router.Group("/category")
    group.POST("/list", auth.Permission("read").Wrap(httpx.HandlerFunc(s.List)))
    group.POST("/create", auth.Permission("create").Wrap(httpx.HandlerFunc(s.Create)))
    group.POST("/update", auth.Permission("update").Wrap(httpx.HandlerFunc(s.Update)))
    group.POST("/delete", auth.Permission("delete").Wrap(httpx.HandlerFunc(s.Delete)))
    group.POST("/detail", auth.Permission("read").Wrap(httpx.HandlerFunc(s.Detail)))
    group.POST("/move", auth.Permission("update").Wrap(httpx.HandlerFunc(s.Move)))
}

// List 分类列表
// @x-apifox-folder "文章分类"
// @Security ApiKeyAuth
// @Summary 分类列表
// @Tags    分类
// @Accept  json
// @Produce json
// @Param   body body api.CategoryListParams true "请求参数"
// @Success 200 {object} api.ApiCategoryList
// @Router  /category/list [post]
func (s *CategoryController) List(ctx *gin.Context) error {
    return httpx.BindJSON[api.CategoryListParams](ctx, func(params api.CategoryListParams) (any, error) {
        return s.service.List(ctx, params)
    })
}

// Detail 查看分类
// @x-apifox-folder "文章分类"
// @Security ApiKeyAuth
// @Summary 查看分类
// @Tags    分类
// @Accept  json
// @Produce json
// @Param   body body apitype.ResourceID true "请求参数"
// @Success 200 {object} api.ApiCategory
// @Router  /category/detail [post]
func (s *CategoryController) Detail(ctx *gin.Context) error {
    return httpx.BindJSON[apitype.ResourceID](ctx, func(params apitype.ResourceID) (any, error) {
        return s.service.FindByID(ctx, params.ID)
    })
}

// Create 创建分类
// @x-apifox-folder "文章分类"
// @Security ApiKeyAuth
// @Summary 创建分类
// @Tags    分类
// @Accept  json
// @Produce json
// @Param   body body api.CategoryCreateParams true "请求参数"
// @Success 200 {object} apitype.ApiCreateResult
// @Router  /category/create [post]
func (s *CategoryController) Create(ctx *gin.Context) error {
    return httpx.BindJSON[api.CategoryCreateParams](ctx, func(params api.CategoryCreateParams) (any, error) {
        return s.service.Create(ctx, params)
    })
}

// Update 更新分类
// @x-apifox-folder "文章分类"
// @Security ApiKeyAuth
// @Summary 更新分类
// @Tags    分类
// @Accept  json
// @Produce json
// @Param   body body api.CategoryUpdateParams true "请求参数"
// @Success 200 {object} apitype.ApiEmptyResult
// @Router  /category/update [post]
func (s *CategoryController) Update(ctx *gin.Context) error {
    return httpx.BindJSON[api.CategoryUpdateParams](ctx, func(params api.CategoryUpdateParams) (any, error) {
        return nil, s.service.Update(ctx, params)
    })
}

// Move 移动分类
// @x-apifox-folder "文章分类"
// @Security ApiKeyAuth
// @Summary 移动分类
// @Tags    分类
// @Accept  json
// @Produce json
// @Param   body body api.CategoryMoveParams true "请求参数"
// @Success 200 {object} apitype.ApiEmptyResult
// @Router  /category/move [post]
func (s *CategoryController) Move(ctx *gin.Context) error {
    return httpx.BindJSON[api.CategoryMoveParams](ctx, func(params api.CategoryMoveParams) (any, error) {
        return nil, s.service.Move(ctx, params.ID, params.TargetID)
    })
}

// Delete 删除分类
// @x-apifox-folder "文章分类"
// @Security ApiKeyAuth
// @Summary 删除分类
// @Tags    分类
// @Accept  json
// @Produce json
// @Param   body body apitype.ResourceID true "请求参数"
// @Success 200 {object} apitype.ApiEmptyResult
// @Router  /category/delete [post]
func (s *CategoryController) Delete(ctx *gin.Context) error {
    return httpx.BindJSON[apitype.ResourceID](ctx, func(params apitype.ResourceID) (any, error) {
        return nil, s.service.Delete(ctx, params.ID)
    })
}
