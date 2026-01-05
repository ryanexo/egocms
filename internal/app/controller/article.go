package controller

import (
    `dpcms/internal/app/controller/internal/httpbinding`
    `dpcms/internal/app/dto`
    `dpcms/internal/app/middleware/authz`
    `dpcms/internal/app/service`
    `dpcms/internal/app/util/contextutil`
    
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

// Create 创建文章
// @x-apifox-folder "文章"
// @Summary 创建文章
// @Tags    文章
// @Accept  json
// @Produce json
// @Param   body body dto.ArticleCreateParams true "请求参数"
// @Success 200 {object} swaggertype.CreateResult
// @Router  /article/create [post]
func (s ArticleController) Create(ctx *gin.Context) {
    httpbinding.BindJSON[dto.ArticleCreateParams](ctx, func(params dto.ArticleCreateParams) (any, error) {
        u, err := contextutil.GetAuthorizedUser(ctx)
        if err != nil {
            return nil, err
        }
        return s.srv.Article.Create(ctx, u, params)
    })
}

// Update 更新文章
// @x-apifox-folder "文章"
// @Summary 更新文章
// @Tags    文章
// @Accept  json
// @Produce json
// @Param   body body dto.ArticleUpdateParams true "请求参数"
// @Success 200 {object} swaggertype.EmptyResult
// @Router  /article/update [post]
func (s ArticleController) Update(ctx *gin.Context) {
    httpbinding.BindJSON[dto.ArticleUpdateParams](ctx, func(params dto.ArticleUpdateParams) (any, error) {
        return nil, s.srv.Article.Update(ctx, params)
    })
}

// Delete 删除文章
// @x-apifox-folder "文章"
// @Summary 删除文章
// @Tags    文章
// @Accept  json
// @Produce json
// @Param   body body dto.ResourceID true "请求参数"
// @Success 200 {object} swaggertype.EmptyResult
// @Router  /article/delete [post]
func (s ArticleController) Delete(ctx *gin.Context) {
    httpbinding.BindJSON[dto.ResourceID](ctx, func(params dto.ResourceID) (any, error) {
        return nil, s.srv.Article.Delete(ctx, params.ID)
    })
}

// Detail 查看文章
// @x-apifox-folder "文章"
// @Summary 查看文章
// @Tags    文章
// @Accept  json
// @Produce json
// @Param   body body dto.ResourceID true "请求参数"
// @Success 200 {object} swaggertype.Article
// @Router  /article/detail [post]
func (s ArticleController) Detail(ctx *gin.Context) {
    httpbinding.BindJSON[dto.ResourceID](ctx, func(params dto.ResourceID) (any, error) {
        return s.srv.Article.FindByID(ctx, params.ID)
    })
}
