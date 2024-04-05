package auth

import (
    `gorm.io/gorm`
)

type RolePermission interface {
    Add(...string)
    Remove(...string)
    Contains(string) bool
}

type DefaultRolePermission struct {
    perms map[string]bool
}

func (s *DefaultRolePermission) Add(perm ...string) {
    for _, v := range perm {
        s.perms[v] = true
    }
}

func (s *DefaultRolePermission) Remove(perm ...string) {
    for _, v := range perm {
        delete(s.perms, v)
    }
}

func (s *DefaultRolePermission) Contains(name string) bool {
    return s.perms[name]
}

type DefaultRoleManager struct {
    db    *gorm.DB
    roles map[string]RolePermission
}

func (s *DefaultRoleManager) Add(name string, perm RolePermission) {
}
