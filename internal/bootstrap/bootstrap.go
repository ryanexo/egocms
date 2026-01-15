package bootstrap

import (
    `cms/internal/httpserver`
    `cms/internal/lifecycle`
    
    `github.com/gin-gonic/gin`
    `github.com/google/wire`
)

var BootstrapProvider = wire.NewSet(
    lifecycle.New,
    New,
)

type Bootstrap struct {
    lifecycle *lifecycle.Lifecycle
    http      *httpserver.Launcher
    mw        httpserver.Middleware
    route     httpserver.Route
}

func (s Bootstrap) Start() error {
    s.http.AddMiddleware(s.mw)
    s.http.AddRoutes(s.route)
    
    err := s.lifecycle.WarmUp()
    if err != nil {
        return err
    }
    
    return s.http.Run(gin.Mode() == gin.DebugMode)
}

func New(lc *lifecycle.Lifecycle, http *httpserver.Launcher, mw httpserver.Middleware, route httpserver.Route) Bootstrap {
    return Bootstrap{lifecycle: lc, http: http, mw: mw, route: route}
}
