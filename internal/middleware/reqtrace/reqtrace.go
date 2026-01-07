package reqtrace

import (
    `dpcms/internal/constant`
    
    `github.com/gin-gonic/gin`
    `github.com/google/uuid`
)

type Middleware gin.HandlerFunc

func New() Middleware {
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
