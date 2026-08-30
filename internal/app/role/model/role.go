package model

import (
    "cms/internal/app/role/errno"
    "cms/internal/infra/store/modeltype"
)

const (
    RoleNameLength        = 64
    RoleDescriptionLength = 255
    RoleInheritanceLimit  = 10
)

type Role struct {
    modeltype.Base
    Name        string `gorm:"type:varchar(64)"`
    Description string `gorm:"type:varchar(255)"`
}

func (role *Role) SetName(name string) error {
    if name == "" {
        return errno.ErrNameRequired
    }
    if len(name) > RoleNameLength {
        return errno.ErrInvalidNameLength
    }
    role.Name = name
    return nil
}

func (role *Role) SetDescription(description string) error {
    if description == "" {
        return errno.ErrDescriptionRequired
    }
    if len(description) > RoleDescriptionLength {
        return errno.ErrInvalidDescriptionLength
    }
    role.Description = description
    return nil
}
