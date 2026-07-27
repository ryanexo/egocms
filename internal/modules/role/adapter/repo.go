package adapter

import (
    "context"
    
    `cms/internal/infra/persistence/gorm/dbscope`
    `cms/internal/infra/persistence/gorm/model`
    `cms/internal/infra/persistence/gorm/gquery`
    "cms/internal/pkg/datatype"
    
    "cms/internal/modules/role/contract"
    "cms/internal/modules/role/internal/dto"
    
    "gorm.io/gen"
)

type roleRepo struct {
    query *gquery.Query
}

func NewRoleRepo(query *gquery.Query) contract.RoleRepo {
    return &roleRepo{query}
}

func (r *roleRepo) CloneWithQuery(q *gquery.Query) contract.RoleRepo {
    return NewRoleRepo(q)
}

func (r *roleRepo) Create(ctx context.Context, data *model.Role) error {
    return r.query.Role.WithContext(ctx).Create(data)
}

func (r *roleRepo) Update(ctx context.Context, data *model.Role) (gen.ResultInfo, error) {
    dao := r.query.Role
    return dao.WithContext(ctx).Where(dao.ID.Eq(data.ID.Raw())).Updates(data)
}

func (r *roleRepo) FindByID(ctx context.Context, id datatype.SafeUint64) (*model.Role, error) {
    dao := r.query.Role
    return dao.WithContext(ctx).Where(dao.ID.Eq(id.Raw())).First()
}

func (r *roleRepo) Delete(ctx context.Context, id datatype.SafeUint64) (gen.ResultInfo, error) {
    dao := r.query.Role
    return dao.WithContext(ctx).Where(dao.ID.Eq(id.Raw())).Delete()
}

func (r *roleRepo) List(ctx context.Context, params *dto.RoleListParams) ([]*model.Role, int64, error) {
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
