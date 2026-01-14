package controller

import (
    `cms/internal/app/article/internal/domain`
    `cms/internal/app/article/internal/dto`
    `cms/internal/app/article/service`
    permissionSrv `cms/internal/app/permission/service`
    `cms/internal/httpserver`
    `cms/internal/infra/casbin`
    `cms/internal/middleware/authz`
    `cms/internal/types`
    `cms/internal/util/contextutil`
    `cms/internal/util/httpbinding`
    
    `github.com/gin-gonic/gin`
)

type ArticleController struct {
    articleSrv *service.ArticleService
    permSrv    *permissionSrv.PermissionService
    auth       *authz.Factory
    casbin     *casbin.RoleCasbin
}

func NewArticleController(
    articleSrv *service.ArticleService,
    permSrv *permissionSrv.PermissionService,
    auth *authz.Factory,
    casbin *casbin.RoleCasbin,
) *ArticleController {
    return &ArticleController{
        articleSrv: articleSrv,
        permSrv:    permSrv,
        auth:       auth,
        casbin:     casbin,
    }
}

func (s ArticleController) Setup(router httpserver.Router) {
    acl := s.auth.AccessControl("article")
    
    g := router.Group("/article", acl.Middleware())
    g.POST("/create", s.Create)
    g.POST("/update", s.Update)
    g.POST("/delete", s.Delete)
    g.POST("/detail", s.Detail)
    g.POST("/submit", s.Submit)
    g.POST("/publish", s.Publish)
    g.POST("/reject", s.Reject)
    g.POST("/republish", s.Republish)
    
    acl.WithRouterOption(
        g,
        authz.WithRouterPermission("/create", "create"),
        authz.WithRouterPermission("/update", "update"),
        authz.WithRouterPermission("/delete", "delete"),
        authz.WithRouterPermission("/detail", "read"),
        authz.WithRouterPermission("/submit", "submit"),
        authz.WithRouterPermission("/publish", "publish"),
        authz.WithRouterPermission("/reject", "reject"),
        authz.WithRouterPermission("/republish", "republish"),
    )
}

// Create
// @x-apifox-folder "文章"
// @Security ApiKeyAuth
// @Summary 创建文章
// @Tags    文章
// @Accept  json
// @Produce json
// @Param   body body dto.ArticleCreateParams true "请求参数"
// @Success 200 {object} types.ApiCreateResult
// @Router  /article/create [post]
func (s ArticleController) Create(ctx *gin.Context) {
    httpbinding.BindJSON[dto.ArticleCreateParams](ctx, func(params dto.ArticleCreateParams) (any, error) {
        u, err := contextutil.GetAuthorizedUser(ctx)
        if err != nil {
            return nil, err
        }
        return s.articleSrv.Create(ctx, u, params)
    })
}

// Update
// @x-apifox-folder "文章"
// @Security ApiKeyAuth
// @Summary 更新文章
// @Tags    文章
// @Accept  json
// @Produce json
// @Param   body body dto.ArticleUpdateParams true "请求参数"
// @Success 200 {object} types.ApiEmptyResult
// @Router  /article/update [post]
func (s ArticleController) Update(ctx *gin.Context) {
    httpbinding.BindJSON[dto.ArticleUpdateParams](ctx, func(params dto.ArticleUpdateParams) (any, error) {
        return nil, s.articleSrv.Update(ctx, params)
    })
}

func (s ArticleController) createActor(ctx *gin.Context) (domain.Actor, error) {
    u, err := contextutil.GetAuthorizedUser(ctx)
    if err != nil {
        return domain.Actor{}, err
    }
    canPublishDirect, err := s.casbin.Enforce(u.RoleID.String(), "article", "publish-direct")
    if err != nil {
        return domain.Actor{}, err
    }
    return domain.Actor{
        CanPublishDirect: canPublishDirect,
    }, nil
}

// Submit
// @x-apifox-folder "文章"
// @Security ApiKeyAuth
// @Summary 提交文章
// @Tags    文章
// @Accept  json
// @Produce json
// @Param   body body types.ResourceID true "请求参数"
// @Success 200 {object} types.ApiEmptyResult
// @Router  /article/submit [post]
func (s ArticleController) Submit(ctx *gin.Context) {
    httpbinding.BindJSON[types.ResourceID](ctx, func(params types.ResourceID) (any, error) {
        actor, err := s.createActor(ctx)
        if err != nil {
            return nil, err
        }
        return nil, s.articleSrv.ChangeStatus(ctx, params.ID, func(status *domain.Status) error {
            return status.WithActor(actor).Submit()
        })
    })
}

// Publish
// @x-apifox-folder "文章"
// @Security ApiKeyAuth
// @Summary 发布文章
// @Tags    文章
// @Accept  json
// @Produce json
// @Param   body body types.ResourceID true "请求参数"
// @Success 200 {object} types.ApiEmptyResult
// @Router  /article/publish [post]
func (s ArticleController) Publish(ctx *gin.Context) {
    httpbinding.BindJSON[types.ResourceID](ctx, func(params types.ResourceID) (any, error) {
        actor, err := s.createActor(ctx)
        if err != nil {
            return nil, err
        }
        return nil, s.articleSrv.ChangeStatus(ctx, params.ID, func(status *domain.Status) error {
            return status.WithActor(actor).Publish()
        })
    })
}

// Offline
// @x-apifox-folder "文章"
// @Security ApiKeyAuth
// @Summary 下线文章
// @Tags    文章
// @Accept  json
// @Produce json
// @Param   body body types.ResourceID true "请求参数"
// @Success 200 {object} types.ApiEmptyResult
// @Router  /article/offline [post]
func (s ArticleController) Offline(ctx *gin.Context) {
    httpbinding.BindJSON[types.ResourceID](ctx, func(params types.ResourceID) (any, error) {
        actor, err := s.createActor(ctx)
        if err != nil {
            return nil, err
        }
        return nil, s.articleSrv.ChangeStatus(ctx, params.ID, func(status *domain.Status) error {
            return status.WithActor(actor).Offline()
        })
    })
}

// Reject
// @x-apifox-folder "文章"
// @Security ApiKeyAuth
// @Summary 拒审文章
// @Tags    文章
// @Accept  json
// @Produce json
// @Param   body body types.ResourceID true "请求参数"
// @Success 200 {object} types.ApiEmptyResult
// @Router  /article/reject [post]
func (s ArticleController) Reject(ctx *gin.Context) {
    httpbinding.BindJSON[types.ResourceID](ctx, func(params types.ResourceID) (any, error) {
        actor, err := s.createActor(ctx)
        if err != nil {
            return nil, err
        }
        return nil, s.articleSrv.ChangeStatus(ctx, params.ID, func(status *domain.Status) error {
            return status.WithActor(actor).Reject()
        })
    })
}

// Republish
// @x-apifox-folder "文章"
// @Security ApiKeyAuth
// @Summary 重新提交审核文章
// @Tags    文章
// @Accept  json
// @Produce json
// @Param   body body types.ResourceID true "请求参数"
// @Success 200 {object} types.ApiEmptyResult
// @Router  /article/republish [post]
func (s ArticleController) Republish(ctx *gin.Context) {
    httpbinding.BindJSON[types.ResourceID](ctx, func(params types.ResourceID) (any, error) {
        actor, err := s.createActor(ctx)
        if err != nil {
            return nil, err
        }
        return nil, s.articleSrv.ChangeStatus(ctx, params.ID, func(status *domain.Status) error {
            return status.WithActor(actor).Republish()
        })
    })
}

// Delete
// @x-apifox-folder "文章"
// @Security ApiKeyAuth
// @Summary 删除文章
// @Tags    文章
// @Accept  json
// @Produce json
// @Param   body body types.ResourceID true "请求参数"
// @Success 200 {object} types.ApiEmptyResult
// @Router  /article/delete [post]
func (s ArticleController) Delete(ctx *gin.Context) {
    httpbinding.BindJSON[types.ResourceID](ctx, func(params types.ResourceID) (any, error) {
        return nil, s.articleSrv.Delete(ctx, params.ID)
    })
}

// Detail
// @x-apifox-folder "文章"
// @Security ApiKeyAuth
// @Summary 查看文章
// @Tags    文章
// @Accept  json
// @Produce json
// @Param   body body types.ResourceID true "请求参数"
// @Success 200 {object} dto.ApiArticle
// @Router  /article/detail [post]
func (s ArticleController) Detail(ctx *gin.Context) {
    httpbinding.BindJSON[types.ResourceID](ctx, func(params types.ResourceID) (any, error) {
        return s.articleSrv.FindByID(ctx, params.ID)
    })
}
