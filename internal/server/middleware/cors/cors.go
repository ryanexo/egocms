package cors

import (
    `strings`

    `github.com/gin-gonic/gin`
)

type Options struct {
    ExposeHeaders    []string
    AllowOrigin      []string
    AllowMethods     []string
    AllowHeaders     []string
    AllowCredentials bool
}

func NewCors(opts Options) gin.HandlerFunc {
    var (
        ExposeHeaders    string
        AllowOrigin      string
        AllowMethods     string
        AllowHeaders     string
        AllowCredentials string
    )

    ExposeHeaders = strings.Join(opts.ExposeHeaders, ",")
    AllowOrigin = strings.Join(opts.AllowOrigin, ",")
    AllowMethods = strings.Join(opts.AllowMethods, ",")
    AllowHeaders = strings.Join(opts.AllowHeaders, ",")

    if opts.AllowCredentials {
        AllowCredentials = "true"
    }

    return func(context *gin.Context) {
        context.Header("Access-Control-Allow-Origin", AllowOrigin)
        context.Header("Access-Control-Allow-Methods", AllowMethods)
        context.Header("Access-Control-Allow-Headers", AllowHeaders)
        context.Header("Access-Control-Allow-Credentials", AllowCredentials)
        context.Header("Access-Control-Expose-Headers", ExposeHeaders)

        context.Next()
    }
}
