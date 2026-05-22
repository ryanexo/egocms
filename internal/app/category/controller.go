package category

import (
    "cms/internal/app/category/internal/dto"
    `cms/internal/httpserver`
    "cms/internal/middleware/authz"
    "cms/internal/util/httpbinding"
    "cms/internal/util/types"
    
    _ "cms/internal/util/httpbinding"
    
    "github.com/gin-gonic/gin"
)

type CategoryController struct {
    categorySrv *CategoryService
    auth        *authz.Factory
}

func NewCategoryController(
    categorySrv *CategoryService,
    auth *authz.Factory,
) *CategoryController {
    return &CategoryController{
        categorySrv: categorySrv,
        auth:        auth,
    }
}

func (s CategoryController) Setup(router httpserver.Router) {
    acl := s.auth.AccessControl("category")
    
    g := router.Group("/category", acl.Middleware())
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
// @Security ApiKeyAuth
// @Summary 分类列表
// @Tags    分类
// @Accept  json
// @Produce json
// @Param   body body dto.CategoryListParams true "请求参数"
// @Success 200 {object} dto.ApiCategoryList
// @Router  /category/list [post]
func (s CategoryController) List(ctx *gin.Context) {
    httpbinding.BindJSON[dto.CategoryListParams](ctx, func(params dto.CategoryListParams) (any, error) {
        return s.categorySrv.List(ctx, params)
    })
}

// Detail 查看分类
// @x-apifox-folder "文章分类"
// @Security ApiKeyAuth
// @Summary 查看分类
// @Tags    分类
// @Accept  json
// @Produce json
// @Param   body body types.ResourceID true "请求参数"
// @Success 200 {object} dto.ApiCategory
// @Router  /category/detail [post]
func (s CategoryController) Detail(ctx *gin.Context) {
    httpbinding.BindJSON[types.ResourceID](ctx, func(params types.ResourceID) (any, error) {
        return s.categorySrv.FindByID(ctx, params.ID)
    })
}

// Create 创建分类
// @x-apifox-folder "文章分类"
// @Security ApiKeyAuth
// @Summary 创建分类
// @Tags    分类
// @Accept  json
// @Produce json
// @Param   body body dto.CategoryCreateParams true "请求参数"
// @Success 200 {object} types.ApiCreateResult
// @Router  /category/create [post]
func (s CategoryController) Create(ctx *gin.Context) {
    httpbinding.BindJSON[dto.CategoryCreateParams](ctx, func(params dto.CategoryCreateParams) (any, error) {
        return s.categorySrv.Create(ctx, params)
    })
}

// Move 移动分类
// @x-apifox-folder "文章分类"
// @Security ApiKeyAuth
// @Summary 移动分类
// @Tags    分类
// @Accept  json
// @Produce json
// @Param   body body dto.CategoryMoveParams true "请求参数"
// @Success 200 {object} types.ApiEmptyResult
// @Router  /category/move [post]
func (s CategoryController) Move(ctx *gin.Context) {
    httpbinding.BindJSON[dto.CategoryMoveParams](ctx, func(params dto.CategoryMoveParams) (any, error) {
        return nil, s.categorySrv.Move(ctx, params.ID, params.TargetID)
    })
}

// Update 更新分类
// @x-apifox-folder "文章分类"
// @Security ApiKeyAuth
// @Summary 更新分类
// @Tags    分类
// @Accept  json
// @Produce json
// @Param   body body dto.CategoryUpdateParams true "请求参数"
// @Success 200 {object} types.ApiEmptyResult
// @Router  /category/update [post]
func (s CategoryController) Update(ctx *gin.Context) {
    httpbinding.BindJSON[dto.CategoryUpdateParams](ctx, func(params dto.CategoryUpdateParams) (any, error) {
        err := s.categorySrv.Update(ctx, params)
        return nil, err
    })
}

// Delete 删除分类
// @x-apifox-folder "文章分类"
// @Security ApiKeyAuth
// @Summary 删除分类
// @Tags    分类
// @Accept  json
// @Produce json
// @Param   body body types.ResourceID true "请求参数"
// @Success 200 {object} types.ApiEmptyResult
// @Router  /category/delete [post]
func (s CategoryController) Delete(ctx *gin.Context) {
    httpbinding.BindJSON[types.ResourceID](ctx, func(params types.ResourceID) (any, error) {
        err := s.categorySrv.Delete(ctx, params.ID)
        return nil, err
    })
}
