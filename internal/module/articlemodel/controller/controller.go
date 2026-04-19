package controller

import (
    `cms/internal/domain/articlemodel/internal/dto`
    articleModelSrv `cms/internal/domain/articlemodel/service`
    `cms/internal/httpserver`
    `cms/internal/middleware/authz`
    `cms/internal/util/httpbinding`
    `cms/internal/util/types`
    
    "github.com/gin-gonic/gin"
)

type ArticleModelController struct {
    articleModelSrv *articleModelSrv.ArticleModelService
    auth            *authz.Factory
}

func NewArticleModelController(
    articleModelSrv *articleModelSrv.ArticleModelService,
    auth *authz.Factory,
) *ArticleModelController {
    return &ArticleModelController{articleModelSrv: articleModelSrv, auth: auth}
}

func (s ArticleModelController) Setup(router httpserver.Router) {
    acl := s.auth.AccessControl("article-model")
    
    g := router.Group("/article-model", acl.Middleware())
    g.POST("/create", s.Create)
    g.POST("/delete", s.Delete)
    g.POST("/detail", s.Detail)
    g.POST("/update", s.Update)
    g.POST("/update-schema", s.UpdateSchema)
    
    acl.WithRouterOption(
        g,
        authz.WithRouterPermission("/create", "create"),
        authz.WithRouterPermission("/update-schema", "update"),
        authz.WithRouterPermission("/update", "update"),
        authz.WithRouterPermission("/delete", "delete"),
        authz.WithRouterPermission("/detail", "read"),
    )
}

// Create 创建文章模型
// @x-apifox-folder "文章模型"
// @Security ApiKeyAuth
// @Summary 创建文章模型
// @Tags    文章模型
// @Accept  json
// @Produce json
// @Param   body body dto.ArticleModelCreateParams true "请求参数"
// @Success 200 {object} types.ApiCreateResult
// @Router  /article-model/create [post]
func (s ArticleModelController) Create(ctx *gin.Context) {
    httpbinding.BindJSON[dto.ArticleModelCreateParams](ctx, func(params dto.ArticleModelCreateParams) (any, error) {
        return s.articleModelSrv.CreateModel(ctx, params)
    })
}

// Update 更新文章模型
// @x-apifox-folder "文章模型"
// @Security ApiKeyAuth
// @Summary 更新文章模型
// @Tags    文章模型
// @Accept  json
// @Produce json
// @Param   body body dto.ArticleModelUpdateParams true "请求参数"
// @Success 200 {object} types.ApiEmptyResult
// @Router  /article-model/update [post]
func (s ArticleModelController) Update(ctx *gin.Context) {
    httpbinding.BindJSON[dto.ArticleModelUpdateParams](ctx, func(params dto.ArticleModelUpdateParams) (any, error) {
        return nil, s.articleModelSrv.UpdateModel(ctx, params)
    })
}

// UpdateSchema 更新文章模型schema
// @x-apifox-folder "文章模型"
// @Security ApiKeyAuth
// @Summary 更新文章模型Schema
// @Tags    文章模型
// @Accept  json
// @Produce json
// @Param   body body dto.ArticleModelSchemaUpdateParams true "请求参数"
// @Success 200 {object} types.ApiEmptyResult
// @Router  /article-model/update-schema [post]
func (s ArticleModelController) UpdateSchema(ctx *gin.Context) {
    httpbinding.BindJSON[dto.ArticleModelSchemaUpdateParams](ctx, func(params dto.ArticleModelSchemaUpdateParams) (any, error) {
        return nil, s.articleModelSrv.UpdateModelSchema(ctx, params)
    })
}

// Delete 删除文章模型
// @x-apifox-folder "文章模型"
// @Security ApiKeyAuth
// @Summary 删除文章模型
// @Tags    文章模型
// @Accept  json
// @Produce json
// @Param   body body types.ResourceID true "请求参数"
// @Success 200 {object} types.ApiEmptyResult
// @Router  /article-model/delete [post]
func (s ArticleModelController) Delete(ctx *gin.Context) {
    httpbinding.BindJSON[types.ResourceID](ctx, func(params types.ResourceID) (any, error) {
        return nil, s.articleModelSrv.DeleteModel(ctx, params.ID)
    })
}

// Detail 查看文章模型
// @x-apifox-folder "文章模型"
// @Security ApiKeyAuth
// @Summary 查看文章模型
// @Tags    文章模型
// @Accept  json
// @Produce json
// @Param   body body types.ResourceID true "请求参数"
// @Success 200 {object} dto.ApiArticleModel
// @Router  /article-model/detail [post]
func (s ArticleModelController) Detail(ctx *gin.Context) {
    httpbinding.BindJSON[types.ResourceID](ctx, func(params types.ResourceID) (any, error) {
        return s.articleModelSrv.FindByID(ctx, params.ID)
    })
}

// List 文章模型列表
// @x-apifox-folder "文章模型"
// @Security ApiKeyAuth
// @Summary 查看文章模型列表
// @Tags    文章模型
// @Accept  json
// @Produce json
// @Param   body body dto.ArticleModelListParams true "请求参数"
// @Success 200 {object} dto.ApiArticleModelList
// @Router  /article-model/list [post]
func (s ArticleModelController) List(ctx *gin.Context) {
    httpbinding.BindJSON[dto.ArticleModelListParams](ctx, func(params dto.ArticleModelListParams) (any, error) {
        return s.articleModelSrv.List(ctx, params)
    })
}
