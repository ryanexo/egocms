package provider

import (
	"reflect"

	"cms/internal/app/category"
	appFile "cms/internal/app/file"
	"cms/internal/app/role"
	"cms/internal/app/user"
	"cms/internal/httpx"
	"cms/internal/util/reflectutil"

	"github.com/google/wire"
)

var ControllerProvider = wire.NewSet(
	wire.Struct(new(ControllerSet), "*"),
	NewRouteRegistrar,
	role.NewRoleController,
	category.NewCategoryController,
	user.NewUserController,
	appFile.NewFileController,
)

type ControllerSet struct {
	Role     *role.Controller
	Category *category.CategoryController
	User     *user.Controller
	File     *appFile.Controller
}

var _ httpx.Route = (*ControllerSet)(nil)

func (c ControllerSet) Setup(router httpx.Router) {
	_ = reflectutil.InvokeImplementedStruct[httpx.Route](c, func(_ reflect.Value, i httpx.Route) error {
		i.Setup(router)
		return nil
	})
}

func NewRouteRegistrar(ctlSet *ControllerSet) httpx.Route {
	return ctlSet
}
