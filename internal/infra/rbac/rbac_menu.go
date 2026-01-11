package rbac

import (
    `github.com/casbin/casbin/v3`
    `gorm.io/gorm`
)

type MenuCasbin struct {
    *casbin.Enforcer
}

func NewMenuCasbin(db *gorm.DB) (*MenuCasbin, error) {
    enforcer, err := New(Options{
        DB:        db,
        TableName: "menu_casbin",
        Model:     DefaultModel(),
    })
    if err != nil {
        return nil, err
    }
    return &MenuCasbin{Enforcer: enforcer}, nil
}
