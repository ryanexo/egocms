package controllers

import (
    `dpcms/internal/app/middleware/auth`
    `dpcms/internal/app/services`
    `dpcms/internal/app/services/types/menu`
    `dpcms/internal/erroz`
    `github.com/gin-gonic/gin`
)

type MenuController struct {
    srv *services.Services
}

func NewMenuController(srv *services.Services) MenuController {
    return MenuController{srv}
}

func (c MenuController) setup(engine *gin.Engine) {
    ac, middleware := auth.New(c.srv)
    
    g := engine.Group("/menu", middleware)
    g.POST("/list", c.List)
    g.POST("/create", c.Create)
    g.POST("/update", c.Update)
    g.POST("/delete", c.Delete)
    g.POST("/detail", c.Detail)
    
    ac.SetObjectName("menu").
        AddGroupPermission(g, "/list", "read").
        AddGroupPermission(g, "/detail", "read").
        AddGroupPermission(g, "/create", "create").
        AddGroupPermission(g, "/update", "update").
        AddGroupPermission(g, "/delete", "delete")
}

func (c MenuController) List(ctx *gin.Context) {
    p := &menu.RetrieveListParams{}
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
    p := &menu.DetailParams{}
    if err := ctx.ShouldBind(p); err != nil {
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
    params := menu.DeleteParams{}
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
