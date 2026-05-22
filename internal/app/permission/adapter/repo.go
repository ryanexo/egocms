package adapter

import (
	"cms/internal/pkg/datatype"
	"context"

	"cms/internal/app/permission/contract"
	"cms/internal/infra/persistence/model"
	"cms/internal/infra/persistence/query"

	"gorm.io/gen"
)

type permissionRepo struct {
	persist *query.Query
}

func NewPermissionRepo(persist *query.Query) contract.PermissionRepo {
	return &permissionRepo{persist: persist}
}

func (s permissionRepo) CloneWithQuery(q *query.Query) contract.PermissionRepo {
	return NewPermissionRepo(q)
}

func (s permissionRepo) Create(ctx context.Context, data *model.Permission) error {
	return s.persist.Permission.WithContext(ctx).Create(data)
}

func (s permissionRepo) Update(ctx context.Context, data *model.Permission) (gen.ResultInfo, error) {
	dao := s.persist.Permission
	return dao.WithContext(ctx).Where(dao.ID.Eq(data.ID.Raw())).Updates(data)
}

func (s permissionRepo) Delete(ctx context.Context, id datatype.SafeUint64) (gen.ResultInfo, error) {
	dao := s.persist.Permission
	return dao.WithContext(ctx).Where(dao.ID.Eq(id.Raw())).Delete()
}

func (s permissionRepo) FindByID(ctx context.Context, id datatype.SafeUint64) (*model.Permission, error) {
	dao := s.persist.Permission
	return dao.WithContext(ctx).Where(dao.ID.Eq(id.Raw())).First()
}

func (s permissionRepo) FindByMenuID(ctx context.Context, id datatype.SafeUint64) ([]*model.Permission, error) {
	dao := s.persist.Permission
	return dao.WithContext(ctx).Where(dao.MenuID.Eq(id.Raw())).Find()
}

func (s permissionRepo) FindByResource(ctx context.Context, resource string) ([]*model.Permission, error) {
	dao := s.persist.Permission
	return dao.WithContext(ctx).Where(dao.Resource.Eq(resource)).Find()
}

func (s permissionRepo) FindNoMenuID(ctx context.Context) ([]*model.Permission, error) {
	dao := s.persist.Permission
	return dao.WithContext(ctx).Where(dao.MenuID.IsNull()).Find()
}

func (s permissionRepo) FindAll(ctx context.Context) ([]*model.Permission, error) {
	return s.persist.Permission.WithContext(ctx).Find()
}
