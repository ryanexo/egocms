package httpserver

import (
    "fmt"
    
    "github.com/gin-gonic/gin"
)

type Routes interface {
    SetupRoutes(s *gin.Engine)
}

type Middleware interface {
    SetupMiddleware(s *gin.Engine)
}

type Launcher struct {
    engine *gin.Engine
    config *Config
}

func (launcher Launcher) addRoutes(routes Routes) {
    routes.SetupRoutes(launcher.engine)
}

func (launcher Launcher) addMiddleware(middleware Middleware) {
    middleware.SetupMiddleware(launcher.engine)
}

func (launcher Launcher) Run() error {
    addr := fmt.Sprintf("%s:%d", launcher.config.Host, launcher.config.Port)
    if launcher.config.SSL == nil {
        return launcher.engine.Run(addr)
    }
    return launcher.engine.RunTLS(
        addr,
        launcher.config.SSL.CertFile,
        launcher.config.SSL.KeyFile,
    )
}
