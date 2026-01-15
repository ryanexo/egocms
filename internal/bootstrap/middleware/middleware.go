package middleware

import (
    `reflect`
    
    permChecker `cms/internal/app/permission/adapter`
    corsOptions `cms/internal/app/setting/adapter`
    tokenParser `cms/internal/app/token/adapter`
    `cms/internal/httpserver`
    `cms/internal/middleware/authz`
    `cms/internal/middleware/cors`
    `cms/internal/middleware/log`
    `cms/internal/middleware/recovery`
    `cms/internal/middleware/reqtrace`
    `cms/internal/util/reflectutil`
    
    "github.com/google/wire"
)

var MiddlewareProvider = wire.NewSet(
    NewMiddlewareRegistrar,
    recovery.New,
    log.New,
    cors.New,
    reqtrace.New,
    authz.NewFactory,
    tokenParser.NewTokenParser,
    permChecker.NewPermissionChecker,
    corsOptions.NewCORSOptions,
)

type middlewareSet struct {
    Recovery recovery.Recovery
    ReqTrace reqtrace.RequestTrace
    Logger   log.Log
    CORS     cors.CORS
}

var _ httpserver.Middleware = (*middlewareSet)(nil)

func (m middlewareSet) Setup(registry httpserver.MiddlewareRegistry) {
    _ = reflectutil.InvokeImplementedStruct[httpserver.Middleware](m, func(_ reflect.Value, md httpserver.Middleware) error {
        md.Setup(registry)
        return nil
    })
}

func NewMiddlewareRegistrar(
    recoveryMdl recovery.Recovery,
    reqTraceMdl reqtrace.RequestTrace,
    loggerMdl log.Log,
    corsMdl cors.CORS,
) httpserver.Middleware {
    return middlewareSet{
        Recovery: recoveryMdl,
        ReqTrace: reqTraceMdl,
        Logger:   loggerMdl,
        CORS:     corsMdl,
    }
}
