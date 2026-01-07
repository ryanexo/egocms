package repo

import (
    `context`
    
    `dpcms/internal/app/role/internal/dto`
    `dpcms/internal/infra/persistence/datatype`
    `dpcms/internal/infra/persistence/dbscope`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/infra/persistence/query`
    
    `gorm.io/gen`
)

type RoleRepo struct {
    query *query.Query
}

func NewRoleRepo(query *query.Query) *RoleRepo {
    return &RoleRepo{query}
}

func (r *RoleRepo) Create(ctx context.Context, data *model.Role) error {
    return r.query.Role.WithContext(ctx).Create(data)
}

func (r *RoleRepo) Update(ctx context.Context, data *model.Role) (gen.ResultInfo, error) {
    dao := r.query.Role
    return dao.WithContext(ctx).Where(dao.ID.Eq(data.ID.Raw())).Updates(data)
}

func (r *RoleRepo) FindByID(ctx context.Context, id datatype.SafeUint64) (*model.Role, error) {
    dao := r.query.Role
    return dao.WithContext(ctx).Where(dao.ID.Eq(id.Raw())).First()
}

func (r *RoleRepo) Delete(ctx context.Context, id datatype.SafeUint64) (gen.ResultInfo, error) {
    dao := r.query.Role
    return dao.WithContext(ctx).Where(dao.ID.Eq(id.Raw())).Delete()
}

func (r *RoleRepo) List(ctx context.Context, params *dto.RoleListParams) ([]*model.Role, int64, error) {
    dao := r.query.Role
    q := dao.WithContext(ctx).Scopes(dbscope.Paginate(params.PageNo, params.PageSize))
    
    if params.Name != nil {
        q = q.Where(dao.Name.Like("%" + *params.Name + "%"))
    }
    if params.Description != nil {
        q = q.Where(dao.Description.Like("%" + *params.Description + "%"))
    }
    
    count, err := q.Count()
    if err != nil {
        return nil, 0, err
    }
    roles, err := q.Find()
    if err != nil {
        return nil, 0, err
    }
    return roles, count, nil
}
