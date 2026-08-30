package adapter

import (
	userModel "cms/internal/app/user/model"
	"cms/internal/middleware/authz"
)

type user struct {
	*userModel.User
}

var _ authz.User = (*user)(nil)

func (s user) Role() string {
	return s.User.Role()
}
