package httpserver

import (
    "path/filepath"
    
    `cms/internal/httpserver/validator`
    
    "github.com/gin-gonic/gin"
    "github.com/gin-gonic/gin/binding"
)

func New(cfg Config) (*Launcher, error) {
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
    
    return &Launcher{config: cfg, engine: engine}, nil
}
