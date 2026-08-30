package provider

import (
	"cms/internal/app/category"
	appFile "cms/internal/app/file"
	"cms/internal/app/role"
	"cms/internal/app/user"
	"cms/internal/modules/setting"

	"github.com/google/wire"
)

var ServiceProvider = wire.NewSet(
	role.NewRoleService,
	setting.NewSettingService,
	category.NewCategoryService,
	user.NewUserService,
	role.NewRoleService,
	appFile.NewFileService,
)
