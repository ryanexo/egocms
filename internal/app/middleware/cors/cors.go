package cors

import (
    "github.com/gin-gonic/gin"
)

type CORS gin.HandlerFunc

type Config struct {
    AllowOrigin      string `json:"allowOrigin" yaml:"allowOrigin"`
    AllowMethods     string `json:"allowMethods" yaml:"allowMethods"`
    AllowHeaders     string `json:"allowHeaders" yaml:"allowHeaders"`
    AllowCredentials bool   `json:"allowCredentials" yaml:"allowCredentials"`
    ExposeHeaders    string `json:"exposeHeaders" yaml:"exposeHeaders"`
}

func New(opts *Config) CORS {
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
