package casbin

import (
    `github.com/casbin/casbin/v3`
    `gorm.io/gorm`
)

type RoleCasbin struct {
    *casbin.Enforcer
}

func NewRoleCasbin(db *gorm.DB) (*RoleCasbin, error) {
    enforcer, err := New(Options{
        DB:        db,
        TableName: "role_casbin",
        Model:     DefaultModel(),
    })
    if err != nil {
        return nil, err
    }
    return &RoleCasbin{Enforcer: enforcer}, nil
}
