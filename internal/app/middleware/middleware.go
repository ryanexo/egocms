package middleware

import (
    `reflect`
    
    `dpcms/internal/app/middleware/cors`
    `dpcms/internal/app/middleware/log`
    `dpcms/internal/app/middleware/recovery`
    `dpcms/internal/app/middleware/reqtrace`
    `dpcms/internal/httpserver`
    
    "github.com/gin-gonic/gin"
    "github.com/google/wire"
)

var ProviderSet = wire.NewSet(
    wire.Struct(new(Middleware), "*"),
    NewMiddlewareRegistrar,
    recovery.New,
    log.New,
    cors.New,
    reqtrace.New,
)

type Middleware struct {
    Recovery recovery.Recovery
    ReqTrace reqtrace.ReqTrace
    Logger   log.LoggerMiddleware
    CORS     cors.CORS
}

var _ httpserver.Middleware = (*Middleware)(nil)

func (m Middleware) SetupMiddleware(engine *gin.Engine) {
    val := reflect.ValueOf(m)
    if val.Kind() == reflect.Ptr {
        val = val.Elem()
    }
    
    handler := reflect.TypeOf(gin.HandlerFunc(nil))
    
    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)
        
        if !field.CanConvert(handler) {
            continue
        }
        
        fn := field.Convert(handler).Interface().(gin.HandlerFunc)
        engine.Use(fn)
    }
}

func NewMiddlewareRegistrar(m *Middleware) httpserver.Middleware {
    return m
}
