package middleware

import (
    `reflect`
    
    `dpcms/api/middleware/cors`
    `dpcms/api/middleware/log`
    `dpcms/api/middleware/recovery`
    `dpcms/server`
    "github.com/gin-gonic/gin"
    "github.com/google/wire"
)

var ProviderSet = wire.NewSet(
    wire.Struct(new(Middleware), "*"),
    NewMiddlewareRegistrar,
    recovery.New,
    log.New,
    cors.New,
)

type Middleware struct {
    Recovery recovery.Recovery
    Logger   log.Logger
    CORS     cors.CORS
}

var _ server.Middleware = (*Middleware)(nil)

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

func NewMiddlewareRegistrar(m *Middleware) server.Middleware {
    return m
}
