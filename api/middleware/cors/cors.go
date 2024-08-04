package cors

import (
    "github.com/gin-gonic/gin"
)

type CORS gin.HandlerFunc

type Options struct {
    AllowOrigin      string `json:"allowOrigin"`
    AllowMethods     string `json:"allowMethods"`
    AllowHeaders     string `json:"allowHeaders"`
    AllowCredentials bool   `json:"allowCredentials"`
    ExposeHeaders    string `json:"exposeHeaders"`
}

func New(opts *Options) CORS {
    return func(context *gin.Context) {
        if opts.AllowOrigin != "" {
            context.Header("Access-Control-Allow-Origin", opts.AllowOrigin)
        }
        if opts.AllowMethods != "" {
            context.Header("Access-Control-Allow-Methods", opts.AllowMethods)
        }
        if opts.AllowHeaders != "" {
            context.Header("Access-Control-Allow-Headers", opts.AllowHeaders)
        }
        if opts.AllowCredentials {
            context.Header("Access-Control-Allow-Credentials", "true")
        }
        if opts.ExposeHeaders != "" {
            context.Header("Access-Control-Expose-Headers", opts.ExposeHeaders)
        }
        context.Next()
    }
}
