package controller

import (
    `dpcms/internal/app/controller/internal/httpbinding`
    artDomain `dpcms/internal/app/domain/article`
    `dpcms/internal/app/dto`
    `dpcms/internal/app/middleware/authz`
    `dpcms/internal/app/service`
    `dpcms/internal/app/util/contextutil`
    `dpcms/internal/infra/persistence/model`
    
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

// Update
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

func (s ArticleController) createActor(u *model.User) (artDomain.Actor, error) {
    enforcer := s.srv.RBAC.GetEnforcer()
    canPublishDirect, err := enforcer.Enforce(u.RoleID.String(), "article", "publish-direct")
    if err != nil {
        return artDomain.Actor{}, err
    }
    return artDomain.Actor{
        CanPublishDirect: canPublishDirect,
    }, nil
}

// Submit
// @x-apifox-folder "文章"
// @Summary 提交文章
// @Tags    文章
// @Accept  json
// @Produce json
// @Param   body body dto.ResourceID true "请求参数"
// @Success 200 {object} swaggertype.EmptyResult
// @Router  /article/submit [post]
func (s ArticleController) Submit(ctx *gin.Context) {
    httpbinding.BindJSON[dto.ResourceID](ctx, func(params dto.ResourceID) (any, error) {
        u, err := contextutil.GetAuthorizedUser(ctx)
        if err != nil {
            return nil, err
        }
        actor, err := s.createActor(u)
        if err != nil {
            return nil, err
        }
        return nil, s.srv.Article.ChangeStatus(ctx, params.ID, func(status *artDomain.Status) error {
            return status.WithActor(actor).Submit()
        })
    })
}

// Publish
// @x-apifox-folder "文章"
// @Summary 发布文章
// @Tags    文章
// @Accept  json
// @Produce json
// @Param   body body dto.ResourceID true "请求参数"
// @Success 200 {object} swaggertype.EmptyResult
// @Router  /article/publish [post]
func (s ArticleController) Publish(ctx *gin.Context) {
    httpbinding.BindJSON[dto.ResourceID](ctx, func(params dto.ResourceID) (any, error) {
        u, err := contextutil.GetAuthorizedUser(ctx)
        if err != nil {
            return nil, err
        }
        actor, err := s.createActor(u)
        if err != nil {
            return nil, err
        }
        return nil, s.srv.Article.ChangeStatus(ctx, params.ID, func(status *artDomain.Status) error {
            return status.WithActor(actor).Publish()
        })
    })
}

// Offline
// @x-apifox-folder "文章"
// @Summary 下线文章
// @Tags    文章
// @Accept  json
// @Produce json
// @Param   body body dto.ResourceID true "请求参数"
// @Success 200 {object} swaggertype.EmptyResult
// @Router  /article/offline [post]
func (s ArticleController) Offline(ctx *gin.Context) {
    httpbinding.BindJSON[dto.ResourceID](ctx, func(params dto.ResourceID) (any, error) {
        u, err := contextutil.GetAuthorizedUser(ctx)
        if err != nil {
            return nil, err
        }
        actor, err := s.createActor(u)
        if err != nil {
            return nil, err
        }
        return nil, s.srv.Article.ChangeStatus(ctx, params.ID, func(status *artDomain.Status) error {
            return status.WithActor(actor).Offline()
        })
    })
}

// Reject
// @x-apifox-folder "文章"
// @Summary 拒审文章
// @Tags    文章
// @Accept  json
// @Produce json
// @Param   body body dto.ResourceID true "请求参数"
// @Success 200 {object} swaggertype.EmptyResult
// @Router  /article/reject [post]
func (s ArticleController) Reject(ctx *gin.Context) {
    httpbinding.BindJSON[dto.ResourceID](ctx, func(params dto.ResourceID) (any, error) {
        u, err := contextutil.GetAuthorizedUser(ctx)
        if err != nil {
            return nil, err
        }
        actor, err := s.createActor(u)
        if err != nil {
            return nil, err
        }
        return nil, s.srv.Article.ChangeStatus(ctx, params.ID, func(status *artDomain.Status) error {
            return status.WithActor(actor).Reject()
        })
    })
}

// Republish
// @x-apifox-folder "文章"
// @Summary 重新提交审核文章
// @Tags    文章
// @Accept  json
// @Produce json
// @Param   body body dto.ResourceID true "请求参数"
// @Success 200 {object} swaggertype.EmptyResult
// @Router  /article/republish [post]
func (s ArticleController) Republish(ctx *gin.Context) {
    httpbinding.BindJSON[dto.ResourceID](ctx, func(params dto.ResourceID) (any, error) {
        u, err := contextutil.GetAuthorizedUser(ctx)
        if err != nil {
            return nil, err
        }
        actor, err := s.createActor(u)
        if err != nil {
            return nil, err
        }
        return nil, s.srv.Article.ChangeStatus(ctx, params.ID, func(status *artDomain.Status) error {
            return status.WithActor(actor).Republish()
        })
    })
}

// Delete
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

// Detail
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
