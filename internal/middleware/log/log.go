package log

import (
    "time"
    
    `cms/internal/constant`
    `cms/internal/httpserver`
    `cms/internal/pkg/logger`
    
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type Log struct {
    logger *logger.Logger
}

func (s Log) Setup(registry httpserver.MiddlewareRegistry) {
    registry.Use(s.Middleware())
}

func (s Log) Middleware() gin.HandlerFunc {
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
        
        traceID, ok := context.Get(constant.RequestTraceIdKey)
        if ok {
            fields = append(fields, zap.Any("trace-id", traceID))
        }
        
        if len(context.Errors) > 0 {
            lvl = zap.ErrorLevel
            fields = append(fields, zap.String("error", context.Errors.String()))
        }
        
        s.logger.Access.Log(lvl, context.Request.URL.Path, fields...)
    }
}

func New(l *logger.Logger) Log {
    return Log{logger: l}
}
