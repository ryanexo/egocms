package recovery

import (
    `errors`
    `net`
    `net/http`
    `net/http/httputil`
    `os`
    `strings`
    
    `cms/internal/erroz`
    `cms/internal/httpserver`
    `cms/internal/infra/logger`
    
    "github.com/gin-gonic/gin"
    `go.uber.org/zap`
)

type Recovery struct {
    logger *logger.Logger
}

func (s Recovery) Setup(registry httpserver.MiddlewareRegistry) {
    registry.Use(s.middleware())
}

func (s Recovery) middleware() gin.HandlerFunc {
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
            
            s.logger.Access.Error(message, fields...)
            s.logger.App.Error(message, fields...)
        }()
        
        ctx.Next()
    }
}

func New(l *logger.Logger) Recovery {
    return Recovery{logger: l}
}
