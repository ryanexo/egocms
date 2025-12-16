package log

import (
    "time"
    
    `dpcms/internal/app/constant`
    
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type Logger gin.HandlerFunc

func New(logger *zap.Logger) Logger {
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
        
        traceId, ok := context.Get(constant.RequestTraceIdKey)
        if ok {
            fields = append(fields, zap.Any("trace-id", traceId))
        }
        
        if len(context.Errors) > 0 {
            lvl = zap.ErrorLevel
            fields = append(fields, zap.String("error", context.Errors.String()))
        }
        
        logger.Log(lvl, context.Request.URL.Path, fields...)
    }
}
