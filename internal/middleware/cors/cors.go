package cors

import (
    `cms/internal/httpserver`
    
    `github.com/gin-gonic/gin`
)

type CORS struct {
    opts Options
}

func (s CORS) Setup(registry httpserver.MiddlewareRegistry) {
    registry.Use(s.middleware())
}

func (s CORS) middleware() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        allowOrigin := s.opts.GetAllowOrigin(ctx)
        allowMethods := s.opts.GetAllowMethods(ctx)
        allowHeaders := s.opts.GetAllowHeaders(ctx)
        allowCredentials := s.opts.GetAllowCredentials(ctx)
        exposeHeaders := s.opts.GetExposeHeaders(ctx)
        
        if allowOrigin != "" {
            ctx.Header("Access-Control-Allow-Origin", allowOrigin)
        }
        if allowMethods != "" {
            ctx.Header("Access-Control-Allow-Methods", allowMethods)
        }
        if allowHeaders != "" {
            ctx.Header("Access-Control-Allow-Headers", allowHeaders)
        }
        if allowCredentials == "true" {
            ctx.Header("Access-Control-Allow-Credentials", "true")
        }
        if exposeHeaders != "" {
            ctx.Header("Access-Control-Expose-Headers", exposeHeaders)
        }
        
        ctx.Next()
    }
}

func New(opts Options) CORS {
    return CORS{opts: opts}
}
