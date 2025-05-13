package server

import (
    "path/filepath"
    
    "github.com/gin-gonic/gin"
    "github.com/gin-gonic/gin/binding"
    "github.com/google/wire"
)

var ProviderSet = wire.NewSet(
    New,
)

func New(cfg *Config, middleware Middleware, routes Routes, validator binding.StructValidator) (*Launcher, error) {
    if cfg.Debug {
        gin.SetMode(gin.DebugMode)
    }
    
    binding.Validator = validator
    
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
