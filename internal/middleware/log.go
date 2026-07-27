package middleware

import (
    "time"
    
    `cms/internal/infra/logger`
    
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type RequestLog gin.HandlerFunc

func NewRequestLog(log *logger.Logger) RequestLog {
    return func(context *gin.Context) {
        startTime := time.Now()
        context.Next()
        cost := time.Since(startTime)
        
        lvl := zap.InfoLevel
        fields := []zap.Field{
            zap.Int("status", context.Writer.Status()),
            zap.String("method", context.Request.Method),
            zap.String("gquery", context.Request.URL.RawQuery),
            zap.String("ip", context.RemoteIP()),
            zap.String("ua", context.Request.UserAgent()),
            zap.Duration("cost", cost),
            zap.String("trace-id", context.GetHeader("x-trace-id")),
        }
        
        if len(context.Errors) > 0 {
            lvl = zap.ErrorLevel
            fields = append(fields, zap.String("error", context.Errors.String()))
        }
        
        log.Access.Log(lvl, context.Request.URL.Path, fields...)
    }
}
