package provider

import (
    `context`
    
    `cms/internal/infra/persistence/contract`
    `cms/internal/infra/persistence/datatype`
    `cms/internal/infra/persistence/model`
    `cms/internal/infra/persistence/query`
)

type fileRepo struct {
    persist *query.Query
}

var _ contract.FileRepo = (*fileRepo)(nil)

func NewFileRepo(persist *query.Query) contract.FileRepo {
    return fileRepo{persist}
}

func (s fileRepo) CloneWithQuery(q *query.Query) contract.FileRepo {
    return NewFileRepo(q)
}

func (s fileRepo) Create(ctx context.Context, data *model.File) error {
    return s.persist.File.WithContext(ctx).Create(data)
}

func (s fileRepo) Delete(ctx context.Context, id datatype.SafeUint64) (gen.ResultInfo, error) {
    dao := s.persist.File
    return dao.WithContext(ctx).Where(dao.ID.Eq(id.Raw())).Delete()
}

func (s fileRepo) FindByID(ctx context.Context, id datatype.SafeUint64) (*model.File, error) {
    dao := s.persist.File
    return dao.WithContext(ctx).Where(dao.ID.Eq(id.Raw())).First()
}

func (s fileRepo) FindByPath(ctx context.Context, path string) (*model.File, error) {
    dao := s.persist.File
    return dao.WithContext(ctx).Where(dao.Path.Eq(path)).First()
}
