package server

import (
    `fmt`

    `GoBlog/internal/server/biz/controller`
    `GoBlog/internal/server/config`
    `GoBlog/internal/server/middleware`
    `github.com/gin-gonic/gin`
)

type Launcher interface {
    Listen() error
}
type listenFunc func() error

func (listen listenFunc) Listen() error {
    return listen()
}

func RegisterHandler(
    cfg *config.Config,
    engine *gin.Engine,
    mdw middleware.Middleware,
    userCtrl controller.UserController,
) Launcher {
    mdw.RegisterForServer(engine)
    userCtrl.RegisterForServer(engine)

    return listenFunc(func() error {
        addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
        if cfg.Server.SSL.IsEnabled() {
            return engine.RunTLS(
                addr,
                cfg.Server.SSL.CertFile,
                cfg.Server.SSL.KeyFile,
            )
        }
        return engine.Run(addr)
    })
}
