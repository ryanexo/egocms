package controller

import (
    `dpcms/internal/app/erroz`
    `dpcms/internal/app/middleware/authz`
    `dpcms/internal/app/service`
    `dpcms/internal/app/service/types`
    `dpcms/internal/app/service/types/menu`
    
    `github.com/gin-gonic/gin`
)

type MenuController struct {
    srv *service.Services
}

func NewMenuController(srv *service.Services) MenuController {
    return MenuController{srv}
}

func (c MenuController) setup(engine *gin.Engine) {
    acl := authz.NewWithRBAC(c.srv, "menu")
    
    g := engine.Group("/menu", acl.Middleware())
    g.POST("/list", c.List)
    g.POST("/create", c.Create)
    g.POST("/update", c.Update)
    g.POST("/delete", c.Delete)
    g.POST("/detail", c.Detail)
    
    acl.WithRouterOption(
        g,
        authz.WithRouterPermission("/list", "read"),
        authz.WithRouterPermission("/detail", "read"),
        authz.WithRouterPermission("/create", "create"),
        authz.WithRouterPermission("/update", "update"),
        authz.WithRouterPermission("/delete", "delete"),
    )
}

func (c MenuController) List(ctx *gin.Context) {
    p := &menu.ListParams{}
    if err := ctx.ShouldBind(p); err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    result, err := c.srv.Menu.List(ctx, p)
    if err != nil {
        erroz.ResolveWithWrite(ctx, err)
    } else {
        erroz.OK.WithOption(erroz.WithData(result)).Write(ctx)
    }
}

func (c MenuController) Detail(ctx *gin.Context) {
    p := types.QueryByIdParam{}
    if err := ctx.ShouldBind(&p); err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    result, err := c.srv.Menu.FindByID(ctx, p.ID)
    if err != nil {
        erroz.ResolveWithWrite(ctx, err)
    } else {
        erroz.OK.WithOption(erroz.WithData(result)).Write(ctx)
    }
}

func (c MenuController) Create(ctx *gin.Context) {
    params := menu.CreateParams{}
    if err := ctx.ShouldBindJSON(&params); err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    result, err := c.srv.Menu.Create(ctx, params)
    if err != nil {
        erroz.ResolveWithWrite(ctx, err)
    } else {
        erroz.OK.WithOption(erroz.WithData(result)).Write(ctx)
    }
}

func (c MenuController) Update(ctx *gin.Context) {
    params := menu.UpdateParams{}
    if err := ctx.ShouldBindJSON(&params); err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    result, err := c.srv.Menu.Update(ctx, params)
    if err != nil {
        erroz.ResolveWithWrite(ctx, err)
    } else {
        erroz.OK.WithOption(erroz.WithData(result)).Write(ctx)
    }
}

func (c MenuController) Delete(ctx *gin.Context) {
    params := types.QueryByIdParam{}
    if err := ctx.ShouldBindJSON(&params); err != nil {
        erroz.ResolveWithWrite(ctx, err)
        return
    }
    err := c.srv.Menu.Delete(ctx, params.ID)
    if err != nil {
        erroz.ResolveWithWrite(ctx, err)
    } else {
        erroz.OK.Write(ctx)
    }
}
