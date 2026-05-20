package bootstrap

import (
    `cms/internal/bootstrap/internal/provider`
    `cms/internal/httpserver`
    
    `github.com/gin-gonic/gin`
    `github.com/google/wire`
)

var BootstrapProvider = wire.NewSet(
    provider.ConfigProvider,
    provider.ServiceProvider,
    provider.InfraProvider,
    provider.RepoProvider,
    provider.ControllerProvider,
    provider.MiddlewareProvider,
    provider.FileDriverProvider,
    New,
)

type Bootstrap struct {
    http  *httpserver.Launcher
    mw    httpserver.Middleware
    route httpserver.Route
}

func (s Bootstrap) Start() error {
    s.http.AddMiddleware(s.mw)
    s.http.AddRoutes(s.route)
    
    return s.http.Run(gin.Mode() == gin.DebugMode)
}

func New(http *httpserver.Launcher, mw httpserver.Middleware, route httpserver.Route) Bootstrap {
    return Bootstrap{http: http, mw: mw, route: route}
}
