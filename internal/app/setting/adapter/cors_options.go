package adapter

import (
    `context`
    
    `cms/internal/app/setting/service`
    `cms/internal/middleware`
)

type CORSOptions struct {
    settingSrv *service.SettingService
}

var _ middleware.Options = (*CORSOptions)(nil)

func (s CORSOptions) GetAllowOrigin(ctx context.Context) string {
    value, _ := s.settingSrv.Get(ctx, "allow-origin")
    return value
}

func (s CORSOptions) GetAllowMethods(ctx context.Context) string {
    value, _ := s.settingSrv.Get(ctx, "allow-methods")
    return value
}

func (s CORSOptions) GetAllowHeaders(ctx context.Context) string {
    value, _ := s.settingSrv.Get(ctx, "allow-headers")
    return value
}

func (s CORSOptions) GetAllowCredentials(ctx context.Context) string {
    value, _ := s.settingSrv.Get(ctx, "allow-origin")
    return value
}

func (s CORSOptions) GetExposeHeaders(ctx context.Context) string {
    value, _ := s.settingSrv.Get(ctx, "expose-headers")
    return value
}

func NewCORSOptions(settingSrv *service.SettingService) middleware.Options {
    return CORSOptions{settingSrv: settingSrv}
}
