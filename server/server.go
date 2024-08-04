package server

import (
    "path/filepath"
    
    `dpcms/packages/validate`
    "github.com/gin-gonic/gin"
    "github.com/gin-gonic/gin/binding"
    "github.com/google/wire"
)

var ServerProviderSet = wire.NewSet(
    New,
)

func New(cfg *Config, middleware Middleware, routes Routes) (*Launcher, error) {
    if cfg.Debug {
        gin.SetMode(gin.DebugMode)
    }
    
    v, err := validate.New()
    if err != nil {
        return nil, err
    }
    binding.Validator = v
    
    engine := gin.New()
    engine.Static(filepath.Base(cfg.StaticDir), filepath.Dir(cfg.StaticDir))
    if cfg.Trust != nil {
        if err := engine.SetTrustedProxies(cfg.Trust); err != nil {
            return nil, err
        }
    }
    
    launcher := &Launcher{config: cfg, engine: engine}
    launcher.addMiddleware(middleware)
    launcher.addRoutes(routes)
    return launcher, nil
}
