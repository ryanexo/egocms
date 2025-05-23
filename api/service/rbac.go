package service

import (
    `fmt`
    
    `dpcms/api/infra`
    `dpcms/enum`
    `github.com/casbin/casbin/v2`
    gormadapter `github.com/casbin/gorm-adapter/v3`
)

type RBACService = *casbin.Enforcer

func NewRBACService(infra *infra.Infra) (RBACService, error) {
    var model = `
    [request_definition]
    r = sub, obj, act
    
    [policy_definition]
    p = sub, obj, act
    
    [role_definition]
    g = _, _
    
    [policy_effect]
    e = some(where (p.eft == allow))
    
    [matchers]
    m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act || r.sub == "%s"
    `
    model = fmt.Sprintf(model, enum.SuperAdminRole)
    adapter, err := gormadapter.NewAdapterByDB(infra.DB)
    if err != nil {
        return nil, err
    }
    
    enforcer, err := casbin.NewEnforcer(model, adapter)
    if err != nil {
        return nil, err
    }
    
    err = enforcer.LoadPolicy()
    if err != nil {
        return nil, err
    }
    
    return enforcer, nil
}
