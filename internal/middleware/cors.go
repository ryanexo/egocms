package middleware

import (
    `context`
    
    `github.com/gin-gonic/gin`
)

type Options interface {
    GetAllowOrigin(context.Context) string
    GetAllowMethods(context.Context) string
    GetAllowHeaders(context.Context) string
    GetAllowCredentials(context.Context) string
    GetExposeHeaders(context.Context) string
}

type CORS gin.HandlerFunc

func NewCORS(opts Options) CORS {
    return func(ctx *gin.Context) {
        allowOrigin := opts.GetAllowOrigin(ctx)
        allowMethods := opts.GetAllowMethods(ctx)
        allowHeaders := opts.GetAllowHeaders(ctx)
        allowCredentials := opts.GetAllowCredentials(ctx)
        exposeHeaders := opts.GetExposeHeaders(ctx)
        
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
