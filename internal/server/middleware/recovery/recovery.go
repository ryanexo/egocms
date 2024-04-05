package recovery

import (
    `errors`
    `net`
    `net/http`
    `net/http/httputil`
    `os`
    `runtime/debug`
    `strings`

    `github.com/gin-gonic/gin`
    `go.uber.org/zap`
)

func New() gin.HandlerFunc {
    return func(context *gin.Context) {
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

                        _ = context.Error(se)
                        context.Abort()
                    }
                }
            }

            stack := string(debug.Stack())
            httpRequest, _ := httputil.DumpRequest(context.Request, false)
            headers := strings.Split(string(httpRequest), "\r\n")

            for idx, header := range headers {
                current := strings.Split(header, ":")

                if current[0] == "Authorization" {
                    headers[idx] = current[0] + ": *"
                }
            }
            gin.Recovery()
            headerString := strings.Join(headers, "\r\n")

            var (
                message string
                fields  = []zap.Field{
                    zap.Any("panic", panicMsg),
                    zap.String("header", headerString),
                }
            )

            if brokenPipe {
                message = "broken pipe"
                context.Abort()
            } else {
                message = "panic"
                fields = append(fields, zap.String("stack", stack))
                context.AbortWithStatus(http.StatusInternalServerError)
            }

            zap.L().Error(message, fields...)
        }()

        context.Next()
    }
}
