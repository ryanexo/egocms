package adapter

import (
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/middleware/authz`
)

type user struct {
    *model.User
}

var _ authz.User = (*user)(nil)

func (s user) Role() string {
    return s.RoleID.String()
}
