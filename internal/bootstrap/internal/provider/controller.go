package provider

import (
    "reflect"
    
    "cms/internal/httpx"
    "cms/internal/util/reflectutil"
    
    "github.com/google/wire"
)

var ControllerProvider = wire.NewSet(
    wire.Struct(new(ControllerSet), "*"),
    NewRouteRegistrar,
)

type ControllerSet struct{}

var _ httpx.Route = (*ControllerSet)(nil)

func (c ControllerSet) Setup(router httpx.Router) {
    _ = reflectutil.InvokeImplementedStruct[httpx.Route](c, func(_ reflect.Value, i httpx.Route) error {
        i.Setup(router)
        return nil
    })
}

func NewRouteRegistrar(ctlSet *ControllerSet) httpx.Route {
    return ctlSet
}
