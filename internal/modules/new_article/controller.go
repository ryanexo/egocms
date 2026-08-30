package new_article

import (
    "context"
    
    "cms/internal/httpx"
    "cms/internal/infra/casbin"
    "cms/internal/middleware/authz"
    `cms/internal/public/apitype`
    `cms/internal/public/jsontype`
    
    "github.com/gin-gonic/gin"
)

type ArticleController struct {
    service    *Service
    auth       *authz.Factory
    roleCasbin *casbin.RoleCasbin
}

func NewArticleController(service *Service, auth *authz.Factory, roleCasbin *casbin.RoleCasbin) *ArticleController {
    return &ArticleController{service: service, auth: auth, roleCasbin: roleCasbin}
}

func (controller *ArticleController) Setup(router httpx.Router) {
    auth := controller.auth.auth("article")
    group := router.Group("/article", auth.Middleware())
    group.POST("/create", httpx.HandlerFunc(controller.CreateDraft))
    group.POST("/update", httpx.HandlerFunc(controller.UpdateDraft))
    group.POST("/submit", httpx.HandlerFunc(controller.Submit))
    group.POST("/publish", httpx.HandlerFunc(controller.Publish))
    group.POST("/reject", httpx.HandlerFunc(controller.Reject))
    group.POST("/offline", httpx.HandlerFunc(controller.Offline))
    group.POST("/republish", httpx.HandlerFunc(controller.Republish))
    group.POST("/delete", httpx.HandlerFunc(controller.Delete))
    group.POST("/detail", httpx.HandlerFunc(controller.Detail))
    
    auth.WithRouterOption(
        group,
        authz.WithRouterPermission("/create", "create"),
        authz.WithRouterPermission("/update", "update"),
        authz.WithRouterPermission("/submit", "submit"),
        authz.WithRouterPermission("/publish", "publish"),
        authz.WithRouterPermission("/reject", "reject"),
        authz.WithRouterPermission("/offline", "offline"),
        authz.WithRouterPermission("/republish", "republish"),
        authz.WithRouterPermission("/delete", "delete"),
        authz.WithRouterPermission("/detail", "read"),
    )
}

func (controller *ArticleController) CreateDraft(ctx *gin.Context) error {
    return httpx.BindJSON[CreateParams](ctx, func(params CreateParams) (any, error) {
        currentUser, err := authz.GetCurrentUser(ctx)
        if err != nil {
            return nil, err
        }
        return controller.service.CreateDraft(ctx, currentUser.UserID(), params)
    })
}

func (controller *ArticleController) UpdateDraft(ctx *gin.Context) error {
    return httpx.BindJSON[UpdateParams](ctx, func(params UpdateParams) (any, error) {
        return nil, controller.service.UpdateDraft(ctx, params)
    })
}

func (controller *ArticleController) Submit(ctx *gin.Context) error {
    return controller.changeWithDirectPermission(ctx, controller.service.Submit)
}

func (controller *ArticleController) Publish(ctx *gin.Context) error {
    return controller.change(ctx, controller.service.Publish)
}

func (controller *ArticleController) Reject(ctx *gin.Context) error {
    return controller.change(ctx, controller.service.Reject)
}

func (controller *ArticleController) Offline(ctx *gin.Context) error {
    return controller.change(ctx, controller.service.Offline)
}

func (controller *ArticleController) Republish(ctx *gin.Context) error {
    return controller.changeWithDirectPermission(ctx, controller.service.Republish)
}

func (controller *ArticleController) Delete(ctx *gin.Context) error {
    return controller.change(ctx, controller.service.Delete)
}

func (controller *ArticleController) Detail(ctx *gin.Context) error {
    return httpx.BindJSON[apitype.ResourceID](ctx, func(params apitype.ResourceID) (any, error) {
        return controller.service.FindByID(ctx, params.ID)
    })
}

func (controller *ArticleController) change(
    ctx *gin.Context,
    operation func(context.Context, jsontype.SafeUint64) error,
) error {
    return httpx.BindJSON[apitype.ResourceID](ctx, func(params apitype.ResourceID) (any, error) {
        return nil, operation(ctx, params.ID)
    })
}

func (controller *ArticleController) changeWithDirectPermission(
    ctx *gin.Context,
    operation func(context.Context, jsontype.SafeUint64, bool) error,
) error {
    return httpx.BindJSON[apitype.ResourceID](ctx, func(params apitype.ResourceID) (any, error) {
        currentUser, err := authz.GetCurrentUser(ctx)
        if err != nil {
            return nil, err
        }
        canPublishDirect, err := controller.roleCasbin.Enforce(currentUser.Role(), "article", "publish-direct")
        if err != nil {
            return nil, err
        }
        return nil, operation(ctx, params.ID, canPublishDirect)
    })
}
