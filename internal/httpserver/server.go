package httpserver

import (
    "path/filepath"
    
    `cms/internal/httpserver/validator`
    
    "github.com/gin-gonic/gin"
    "github.com/gin-gonic/gin/binding"
)

func New(cfg Config, middleware Middleware, routes Route) (*Launcher, error) {
    if cfg.Debug {
        gin.SetMode(gin.DebugMode)
    }
    
    binding.Validator = validator.New()
    
    engine := gin.New()
    
    if cfg.StaticDir != "" {
        engine.Static(filepath.Base(cfg.StaticDir), filepath.Dir(cfg.StaticDir))
    }
    if cfg.Trust != nil {
        if err := engine.SetTrustedProxies(cfg.Trust); err != nil {
            return nil, err
        }
    }
    if cfg.MaxMemory == 0 {
        engine.MaxMultipartMemory = 8 << 20
    }
    
    launcher := &Launcher{config: cfg, engine: engine}
    launcher.addMiddleware(middleware)
    launcher.addRoutes(routes)
    return launcher, nil
}
