package provider

import (
    `reflect`
    
    permChecker `cms/internal/app/permission/adapter`
    corsOptions `cms/internal/app/setting/adapter`
    tokenParser `cms/internal/app/token/adapter`
    `cms/internal/httpserver`
    `cms/internal/middleware`
    `cms/internal/middleware/authz`
    
    `github.com/gin-gonic/gin`
    "github.com/google/wire"
)

var MiddlewareProvider = wire.NewSet(
    wire.Struct(new(Middleware), "*"),
    NewMiddlewareRegistrar,
    authz.NewFactory,
    middleware.NewRecovery,
    middleware.NewCORS,
    middleware.NewErrorFallback,
    middleware.NewRequestTrace,
    middleware.NewRequestLog,
    tokenParser.NewTokenParser,
    permChecker.NewPermissionChecker,
    corsOptions.NewCORSOptions,
)

type Middleware struct {
    Recovery   middleware.Recovery
    ReqTrace   middleware.RequestTrace
    RequestLog middleware.RequestLog
    Cors       middleware.CORS
}

var _ httpserver.Middleware = (*Middleware)(nil)

func (m Middleware) Setup(registry httpserver.MiddlewareRegistry) {
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

func NewMiddlewareRegistrar(md Middleware) httpserver.Middleware {
    return md
}
