package controller

import (
    `dpcms/internal/app/controller/internal/httpbinding`
    `dpcms/internal/app/dto`
    `dpcms/internal/app/middleware/authz`
    `dpcms/internal/app/service`
    
    `github.com/gin-gonic/gin`
)

type ArticleController struct {
    srv *service.Services
}

func NewArticleController(srv *service.Services) ArticleController {
    return ArticleController{srv: srv}
}

func (s ArticleController) setup(engine *gin.Engine) {
    acl := authz.NewWithRBAC(s.srv, "article")
    
    g := engine.Group("/article", acl.Middleware())
    g.POST("/create", s.Create)
    g.POST("/update", s.Update)
    g.POST("/delete", s.Delete)
    g.POST("/detail", s.Detail)
    
    acl.WithRouterOption(
        g,
        authz.WithRouterPermission("/create", "create"),
        authz.WithRouterPermission("/update", "update"),
        authz.WithRouterPermission("/delete", "delete"),
        authz.WithRouterPermission("/detail", "read"),
    )
}

func (s ArticleController) Create(ctx *gin.Context) {
    httpbinding.BindJSON[dto.ArticleCreateParams](ctx, func(params dto.ArticleCreateParams) (any, error) {
        return nil, s.srv.Article.Create(ctx, params)
    })
}

func (s ArticleController) Update(ctx *gin.Context) {
    httpbinding.BindJSON[dto.ArticleUpdateParams](ctx, func(params dto.ArticleUpdateParams) (any, error) {
        return nil, s.srv.Article.Update(ctx, params)
    })
}

func (s ArticleController) Delete(ctx *gin.Context) {
    httpbinding.BindJSON[dto.ResourceID](ctx, func(params dto.ResourceID) (any, error) {
        return nil, s.srv.Article.Delete(ctx, params.ID)
    })
}

func (s ArticleController) Detail(ctx *gin.Context) {
    httpbinding.BindJSON[dto.ResourceID](ctx, func(params dto.ResourceID) (any, error) {
        return s.srv.Article.FindByID(ctx, params.ID)
    })
}
