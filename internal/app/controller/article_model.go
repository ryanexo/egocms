package controller

import (
    `dpcms/internal/app/controller/internal/httpbinding`
    `dpcms/internal/app/dto`
    `dpcms/internal/app/middleware/authz`
    `dpcms/internal/app/service`
    
    `github.com/gin-gonic/gin`
)

type ArticleModelController struct {
    srv *service.Services
}

func NewArticleModelController(srv *service.Services) ArticleModelController {
    return ArticleModelController{srv: srv}
}

func (s ArticleModelController) setup(engine *gin.Engine) {
    acl := authz.NewWithRBAC(s.srv, "article-model")
    
    g := engine.Group("/article-model", acl.Middleware())
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
// @Tags    article
// @Accept  json
// @Produce json
// @Param   body    body   dto.ArticleModelCreateParams    true    "test"
// @Success 200 {null}  nil
// @Router  /article-model/create [post]
func (s ArticleModelController) Create(ctx *gin.Context) {
    httpbinding.BindJSON[dto.ArticleModelCreateParams](ctx, func(params dto.ArticleModelCreateParams) (any, error) {
        return s.srv.ArticleModel.Create(ctx, params)
    })
}

func (s ArticleModelController) Update(ctx *gin.Context) {
    httpbinding.BindJSON[dto.ArticleModelUpdateParams](ctx, func(params dto.ArticleModelUpdateParams) (any, error) {
        return nil, s.srv.ArticleModel.Update(ctx, params)
    })
}

func (s ArticleModelController) UpdateSchema(ctx *gin.Context) {
    httpbinding.BindJSON[dto.ArticleModelSchemaUpdateParams](ctx, func(params dto.ArticleModelSchemaUpdateParams) (any, error) {
        return nil, s.srv.ArticleModel.UpdateSchema(ctx, params)
    })
}

func (s ArticleModelController) Delete(ctx *gin.Context) {
    httpbinding.BindJSON[dto.ResourceID](ctx, func(params dto.ResourceID) (any, error) {
        return nil, s.srv.ArticleModel.DeleteModel(ctx, params.ID)
    })
}

func (s ArticleModelController) Detail(ctx *gin.Context) {
    httpbinding.BindJSON[dto.ResourceID](ctx, func(params dto.ResourceID) (any, error) {
        return s.srv.ArticleModel.FindByID(ctx, params.ID)
    })
}
