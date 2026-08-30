package role

import (
    "context"
    
    "cms/internal/app/role/api"
    "cms/internal/app/role/model"
    "cms/internal/infra/store/gorm/dbscope"
    "cms/internal/infra/store/gorm/gquery"
    
    "gorm.io/gen"
)

type repo struct {
    q *gquery.Query
}

func NewRoleRepo(q *gquery.Query) Repo {
    return &repo{q: q}
}

func (repo *repo) Create(ctx context.Context, role *model.Role) error {
    return repo.q.Role.WithContext(ctx).Create(role)
}

func (repo *repo) Update(ctx context.Context, role *model.Role) (gen.ResultInfo, error) {
    roleQuery := repo.q.Role
    return roleQuery.WithContext(ctx).
        Where(roleQuery.ID.Eq(role.ID)).
        Select(roleQuery.Name, roleQuery.Description).
        Updates(role)
}

func (repo *repo) Delete(ctx context.Context, id uint64) error {
    roleQuery := repo.q.Role
    _, err := roleQuery.WithContext(ctx).Where(roleQuery.ID.Eq(id)).Delete()
    return err
}

func (repo *repo) FindByID(ctx context.Context, id uint64) (*model.Role, error) {
    roleQuery := repo.q.Role
    return roleQuery.WithContext(ctx).Where(roleQuery.ID.Eq(id)).First()
}

func (repo *repo) List(ctx context.Context, params *api.RoleListParams) ([]*model.Role, int64, error) {
    roleQuery := repo.q.Role
    query := roleQuery.WithContext(ctx).Scopes(dbscope.Paginate(params.PageNo, params.PageSize))
    if params.Name != nil {
        query = query.Where(roleQuery.Name.Like("%" + *params.Name + "%"))
    }
    if params.Description != nil {
        query = query.Where(roleQuery.Description.Like("%" + *params.Description + "%"))
    }
    
    total, err := query.Count()
    if err != nil {
        return nil, 0, err
    }
    roles, err := query.Find()
    if err != nil {
        return nil, 0, err
    }
    return roles, total, nil
}
