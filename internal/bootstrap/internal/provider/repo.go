package provider

import (
	"cms/internal/app/category"
	appFile "cms/internal/app/file"
	"cms/internal/app/role"
	"cms/internal/app/setting"
	"cms/internal/app/user"

	"github.com/google/wire"
)

var RepoProvider = wire.NewSet(
	category.NewCategoryRepo,
	role.NewRoleRepo,
	setting.NewSettingRepo,
	user.NewUserRepo,
	appFile.NewFileRepo,
	wire.Bind(new(appFile.Repository), new(*appFile.Repo)),
)
