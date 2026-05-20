package middleware

import (
    `cms/internal/infra/logger`
    
    "github.com/gin-gonic/gin"
    `go.uber.org/zap`
    `go.uber.org/zap/zapio`
)

type Recovery gin.HandlerFunc

func NewRecovery(l *logger.Logger) Recovery {
    return Recovery(
        gin.RecoveryWithWriter(&zapio.Writer{
            Log:   l.App,
            Level: zap.ErrorLevel,
        }),
    )
}
