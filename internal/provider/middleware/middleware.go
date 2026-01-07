package middleware

import (
    `reflect`
    
    permChecker `dpcms/internal/app/permission/auth`
    tokenParser `dpcms/internal/app/token/auth`
    `dpcms/internal/httpserver`
    `dpcms/internal/middleware/authz`
    `dpcms/internal/middleware/cors`
    `dpcms/internal/middleware/log`
    `dpcms/internal/middleware/recovery`
    `dpcms/internal/middleware/reqtrace`
    
    "github.com/gin-gonic/gin"
    "github.com/google/wire"
)

var MiddlewareProvider = wire.NewSet(
    wire.Struct(new(Middleware), "*"),
    NewMiddlewareRegistrar,
    recovery.New,
    log.New,
    cors.New,
    reqtrace.New,
    authz.NewBuilder,
    tokenParser.NewTokenParser,
    permChecker.NewPermissionChecker,
)

type Middleware struct {
    Recovery recovery.Middleware
    ReqTrace reqtrace.Middleware
    Logger   log.Middleware
    CORS     cors.Middleware
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
