package middleware

import (
    `GoBlog/internal/server/middleware/accesslog`
    `GoBlog/internal/server/middleware/recovery`
    `github.com/gin-gonic/gin`
)

type Middleware interface {
    RegisterForServer(*gin.Engine)
}

type middlewareList []gin.HandlerFunc

func (m middlewareList) RegisterForServer(engine *gin.Engine) {
    for _, handler := range m {
        engine.Use(handler)
    }
}

func New() Middleware {
    return middlewareList{
        recovery.New(),
        accesslog.New(),
    }
}
