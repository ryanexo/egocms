package adapter

import (
    `cms/internal/infra/persistence/gorm/model`
    `cms/internal/middleware/authz`
)

type user struct {
    *model.User
}

var _ authz.User = (*user)(nil)

func (s user) Role() string {
    return s.RoleID.String()
}
