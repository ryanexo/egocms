package provider

import (
    `reflect`
    
    `cms/internal/httpx`
    
    `github.com/gin-gonic/gin`
    "github.com/google/wire"
)

var MiddlewareProvider = wire.NewSet(
    wire.Struct(new(Middleware), "*"),
    NewMiddlewareRegistrar,
    // authz.NewFactory,
    // middleware.NewRecovery,
    // middleware.NewCORS,
    // middleware.NewErrorFallback,
    // middleware.NewRequestTrace,
    // middleware.NewRequestLog,
    // permChecker.NewPermissionChecker,
    // corsOptions.NewCORSOptions,
)

type Middleware struct {
    // Recovery   middleware.Recovery
    // ReqTrace   middleware.RequestTrace
    // RequestLog middleware.RequestLog
    // Cors       middleware.CORS
}

var _ httpx.Middleware = (*Middleware)(nil)

func (m Middleware) Setup(registry httpx.MiddlewareRegistry) {
    ref := reflect.ValueOf(m)
    handlerType := reflect.TypeOf((gin.HandlerFunc)(nil))
    
    for i := 0; i < ref.NumField(); i++ {
        field := ref.Field(i)
        
        for {
            if field.Kind() == reflect.Func && field.CanConvert(handlerType) {
                fn := field.Convert(handlerType).Interface().(gin.HandlerFunc)
                registry.Use(fn)
            } else if field.Kind() == reflect.Ptr {
                field = field.Elem()
                continue
            } else {
                break
            }
        }
    }
}

func NewMiddlewareRegistrar(md Middleware) httpx.Middleware {
    return md
}
