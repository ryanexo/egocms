package casbin

import (
    `github.com/casbin/casbin/v3`
    `github.com/casbin/casbin/v3/model`
    gormadapter `github.com/casbin/gorm-adapter/v3`
    `gorm.io/gorm`
)

type Options struct {
    DB        *gorm.DB
    Prefix    string
    TableName string
    Model     string
}

func New(opt Options) (*casbin.Enforcer, error) {
    m, err := model.NewModelFromString(opt.Model)
    if err != nil {
        return nil, err
    }
    
    var adapter *gormadapter.Adapter
    
    if opt.TableName == "" {
        adapter, err = gormadapter.NewAdapterByDB(opt.DB)
    } else {
        adapter, err = gormadapter.NewAdapterByDBUseTableName(opt.DB, opt.Prefix, opt.TableName)
    }
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
    
    return enforcer, nil
}
