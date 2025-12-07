package reqtrace

import (
    `dpcms/internal/app/constant`
    
    `github.com/gin-gonic/gin`
    `github.com/google/uuid`
)

type ReqTrace gin.HandlerFunc

func New() ReqTrace {
    return func(context *gin.Context) {
        id, err := uuid.NewV7()
        if err == nil {
            context.Set(constant.TraceIdKey, id.String())
            context.Header(constant.TraceIdKey, id.String())
        } else {
            _ = context.Error(err)
        }
        context.Next()
    }
}
