package controller

import (
    `dpcms/internal/app/middleware/authz`
    `dpcms/internal/app/service`
    
    `github.com/gin-gonic/gin`
)

type ArticleController struct {
    srv *service.Services
}

func NewArticleController(srv *service.Services) *ArticleController {
    return &ArticleController{srv: srv}
}

func (s *ArticleController) setup(engine *gin.Engine) {
    acl := authz.NewWithRBAC(s.srv, "article")
    
    g := engine.Group("/article", acl.Middleware())
    g.POST("")
}
