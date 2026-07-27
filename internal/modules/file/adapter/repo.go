package adapter

import (
    "context"
    
    `cms/internal/infra/persistence/gorm/model`
    `cms/internal/infra/persistence/gorm/gquery`
    "cms/internal/pkg/datatype"
    
    "cms/internal/modules/file/contract"
    
    "gorm.io/gen"
)

type fileRepo struct {
    persist *gquery.Query
}

var _ contract.FileRepo = (*fileRepo)(nil)

func NewFileRepo(persist *gquery.Query) contract.FileRepo {
    return fileRepo{persist}
}

func (s fileRepo) CloneWithQuery(q *gquery.Query) contract.FileRepo {
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
