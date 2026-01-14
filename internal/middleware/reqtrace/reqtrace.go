package reqtrace

import (
    `cms/internal/constant`
    `cms/internal/httpserver`
    
    `github.com/gin-gonic/gin`
    `github.com/google/uuid`
)

type RequestTrace gin.HandlerFunc

func (s RequestTrace) Setup(registry httpserver.MiddlewareRegistry) {
    registry.Use(gin.HandlerFunc(s))
}

func New() RequestTrace {
    return func(context *gin.Context) {
        id, err := uuid.NewV7()
        if err == nil {
            context.Set(constant.RequestTraceIdKey, id.String())
            context.Header(constant.RequestTraceIdKey, id.String())
        } else {
            _ = context.Error(err)
        }
        context.Next()
    }
}
