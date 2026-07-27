package httpx

import (
    "fmt"
    
    "github.com/gin-gonic/gin"
    swaggerFiles `github.com/swaggo/files`
    ginSwagger `github.com/swaggo/gin-swagger`
)

type Launcher struct {
    engine *gin.Engine
    config Config
}

func (launcher Launcher) AddRoutes(routes Route) {
    routes.Setup(launcher.engine)
}

func (launcher Launcher) AddMiddleware(middleware Middleware) {
    middleware.Setup(launcher.engine)
}

func (launcher Launcher) Run(enableDoc bool) error {
    addr := fmt.Sprintf("%s:%d", launcher.config.Host, launcher.config.Port)
    if enableDoc {
        launcher.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    }
    if launcher.config.SSL == nil || launcher.config.SSL.KeyFile == "" || launcher.config.SSL.CertFile == "" {
        return launcher.engine.Run(addr)
    }
    return launcher.engine.RunTLS(
        addr,
        launcher.config.SSL.CertFile,
        launcher.config.SSL.KeyFile,
    )
}
