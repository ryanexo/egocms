package service

import (
    `fmt`
    
    `dpcms/internal/app/constant`
    `dpcms/internal/infra`
    
    `github.com/casbin/casbin/v2`
    `github.com/casbin/casbin/v2/model`
    gormadapter `github.com/casbin/gorm-adapter/v3`
)

type RBACService struct {
    enforcer *casbin.Enforcer
}

func NewRBACService(infra *infra.Infra) (*RBACService, error) {
    var modelString = `
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
    modelString = fmt.Sprintf(modelString, constant.SuperAdminRole)
    m, err := model.NewModelFromString(modelString)
    if err != nil {
        return nil, err
    }
    
    adapter, err := gormadapter.NewAdapterByDB(infra.DB)
    if err != nil {
        return nil, err
    }
    
    enforcer, err := casbin.NewEnforcer(m, adapter)
    if err != nil {
        return nil, err
    }
    
    err = enforcer.LoadPolicy()
    if err != nil {
        return nil, err
    }
    
    return &RBACService{enforcer: enforcer}, nil
}

func (s *RBACService) GetEnforcer() *casbin.Enforcer {
    return s.enforcer
}
