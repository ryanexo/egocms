package middleware

import (
    `github.com/gin-gonic/gin`
    `github.com/google/uuid`
)

type RequestTrace gin.HandlerFunc

func NewRequestTrace() RequestTrace {
    return func(context *gin.Context) {
        id, err := uuid.NewV7()
        if err == nil {
            context.Header("x-trace-id", id.String())
        } else {
            _ = context.Error(err)
        }
        context.Next()
    }
}
