package adapter

import (
    `context`
    
    `cms/internal/middleware`
    `cms/internal/modules/setting`
)

type CORSOptions struct {
    settingSrv *setting.SettingService
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

func NewCORSOptions(settingSrv *setting.SettingService) middleware.Options {
    return CORSOptions{settingSrv: settingSrv}
}
