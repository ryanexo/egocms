package accesslog

import (
    `time`

    `github.com/gin-gonic/gin`
    `go.uber.org/zap`
)

func New() gin.HandlerFunc {
    return func(context *gin.Context) {
        startTime := time.Now()
        context.Next()
        cost := time.Since(startTime)

        lvl := zap.InfoLevel
        fields := []zap.Field{
            zap.Int("status", context.Writer.Status()),
            zap.String("method", context.Request.Method),
            zap.String("query", context.Request.URL.RawQuery),
            zap.String("ip", context.RemoteIP()),
            zap.String("ua", context.Request.UserAgent()),
            zap.Duration("cost", cost),
        }

        if len(context.Errors) > 0 {
            lvl = zap.ErrorLevel
            for _, e := range context.Errors {
                fields = append(fields, zap.Error(e))
            }
        }

        zap.L().Log(lvl, context.Request.URL.Path, fields...)
    }
}
