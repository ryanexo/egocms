package recovery

import (
    "errors"
    "net"
    "net/http"
    "net/http/httputil"
    "os"
    "strings"
    
    `dpcms/internal/erroz`
    `dpcms/internal/infra/logger`
    
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type Middleware gin.HandlerFunc

func (fn Middleware) Setup(engine *gin.Engine) {
    engine.Use(gin.HandlerFunc(fn))
}

func New(log *logger.Logger) Middleware {
    return func(ctx *gin.Context) {
        if !gin.IsDebugging() {
            _ = os.Stdout.Close()
            _ = os.Stderr.Close()
        }
        
        defer func() {
            panicMsg := recover()
            if panicMsg == nil {
                return
            }
            
            var brokenPipe bool
            if ne, ok := panicMsg.(*net.OpError); ok { // nolint
                var se *os.SyscallError
                
                if errors.As(ne, &se) {
                    seStr := strings.ToLower(se.Error())
                    
                    if strings.Contains(seStr, "broken pipe") || strings.Contains(seStr, "connection reset by peer") {
                        brokenPipe = true
                        
                        _ = ctx.Error(se)
                        ctx.Abort()
                    }
                }
            }
            
            httpRequest, _ := httputil.DumpRequest(ctx.Request, false)
            
            var (
                message string
                fields  = []zap.Field{
                    zap.Any("panic", panicMsg),
                    zap.ByteString("request", httpRequest),
                }
            )
            
            if brokenPipe {
                message = "[broken pipe]"
                ctx.Abort()
            } else {
                message = "panic"
                erroz.Unknown.WithOption(
                    erroz.WithDebug(panicMsg), erroz.WithStatus(http.StatusInternalServerError),
                ).Write(ctx)
            }
            
            log.Access.Error(message, fields...)
            log.App.Error(message, fields...)
        }()
        
        ctx.Next()
    }
}
