package provider

import (
    "context"
    "fmt"
    `strings`
    
    `cms/internal/app/category/domain`
    `cms/internal/app/category/internal/dto`
    `cms/internal/infra/persist/contract`
    `cms/internal/infra/persist/datatype`
    "cms/internal/infra/persist/model"
    "cms/internal/infra/persist/query"
    "cms/internal/infra/persist/scope"
    
    `gorm.io/gen`
)

type categoryRepo struct {
    q *query.Query
}

func NewCategoryRepo(q *query.Query) domain.CategoryRepo {
    return &categoryRepo{q}
}

func (s *categoryRepo) CloneWithQuery(q *query.Query) contract.CategoryRepo {
    return NewCategoryRepo(q)
}

func (s *categoryRepo) Create(ctx context.Context, data domain.Category) error {
    m := model.Category{
        ParentID: data.ParentId(),
        Sequence: data.Sequence(),
        Name:     data.Name(),
        Path:     data.Path().String(),
        Type:     data.Typ(),
        Visible:  0,
        SEO: &model.CategorySeo{
            Title:       data.Seo().Title(),
            Keywords:    strings.Join(data.Seo().Keywords(), ","),
            Description: data.Seo().Description(),
        },
    }
    u, _ := data.URL()
    if u != nil {
        url := u.String()
        m.URL = &url
    }
    return s.q.Category.WithContext(ctx).Create(&m)
}

func (s *categoryRepo) CreateSubtree(ctx context.Context, id datatype.SafeUint64, parentId datatype.SafeUint64) error {
    return s.q.CategoryContext.WithContext(ctx).CreateSubtree(id.Raw(), parentId.Raw())
}

func (s *categoryRepo) Update(ctx context.Context, data domain.Category) (gen.ResultInfo, error) {
    return s.q.Category.WithContext(ctx).Where(
        s.q.Category.ID.Eq(data.Id().Raw()),
    ).Updates(data)
}

func (s *categoryRepo) Move(ctx context.Context, from datatype.SafeUint64, to datatype.SafeUint64) error {
    ctxDao := s.q.CategoryContext
    catDao := s.q.Category
    fromId := from.Raw()
    toId := to.Raw()
    err := ctxDao.WithContext(ctx).UnbindRelationships(fromId)
    if err != nil {
        return err
    }
    err = ctxDao.WithContext(ctx).ReBindRelationships(fromId, toId)
    if err != nil {
        return err
    }
    _, err = catDao.WithContext(ctx).Where(catDao.ID.Eq(fromId)).Update(catDao.ParentID, toId)
    return err
}

func (s *categoryRepo) Delete(ctx context.Context, id datatype.SafeUint64) error {
    ctxDao := s.q.CategoryContext
    catDao := s.q.Category
    seoDao := s.q.CategorySeo
    sqlStr := `DELETE FROM %[1]s WHERE id IN ( SELECT d_id FROM ( SELECT t.%[4]s AS d_id FROM %[2]s AS t WHERE %[3]s = ? ) )`
    deleteCategorySql := fmt.Sprintf(
        sqlStr,
        catDao.TableName(),
        ctxDao.TableName(),
        ctxDao.Ancestor.ColumnName(),
        ctxDao.Descendant.ColumnName(),
    )
    deleteSeoSql := fmt.Sprintf(
        sqlStr,
        seoDao.TableName(),
        ctxDao.TableName(),
        ctxDao.Ancestor.ColumnName(),
        ctxDao.Descendant.ColumnName(),
    )
    
    err := s.q.Category.WithContext(ctx).UnderlyingDB().Exec(deleteCategorySql, id).Error
    if err != nil {
        return err
    }
    
    err = s.q.CategorySeo.WithContext(ctx).UnderlyingDB().Exec(deleteSeoSql, id).Error
    if err != nil {
        return err
    }
    
    err = s.q.CategoryContext.WithContext(ctx).DropSubtree(id.Raw())
    if err != nil {
        return err
    }
    
    return nil
}

func (s *categoryRepo) FindByID(ctx context.Context, id datatype.SafeUint64) (*model.Category, error) {
    return s.q.Category.WithContext(ctx).Preload(s.q.Category.SEO).Where(
        s.q.Category.ID.Eq(id.Raw()),
    ).First()
}

func (s *categoryRepo) FindByIDWithAncestor(ctx context.Context, ancestor datatype.SafeUint64, descendant datatype.SafeUint64) (*model.CategoryContext, error) {
    ctxDao := s.q.CategoryContext
    return ctxDao.WithContext(ctx).Where(
        ctxDao.Ancestor.Eq(ancestor.Raw()),
        ctxDao.Descendant.Eq(descendant.Raw()),
    ).First()
}

func (s *categoryRepo) List(ctx context.Context, params dto.CategoryQuery) ([]*model.Category, int64, error) {
    dao := s.q.Category
    q := dao.WithContext(ctx).Scopes(scope.Paginate(params.PageNo, params.PageSize))
    if params.ID != nil {
        q = q.Where(dao.ID.Eq(params.ID.Raw()))
    }
    if params.ParentID != nil {
        q = q.Where(dao.ParentID.Eq(params.ParentID.Raw()))
    }
    if params.Type != nil {
        q = q.Where(dao.Type.Eq(*params.Type))
    }
    if params.Name != nil {
        q = q.Where(dao.Name.Like("%" + *params.Name + "%"))
    }
    if params.Path != nil {
        q = q.Where(dao.Path.Eq(*params.Path))
    }
    if params.Visible != nil {
        q = q.Where(dao.Visible.Eq(params.Visible.Raw()))
    }
    
    total, err := q.Count()
    if err != nil {
        return nil, 0, err
    }
    
    queryResult, err := q.Preload(dao.SEO).Find()
    if err != nil {
        return nil, 0, err
    }
    
    return queryResult, total, nil
}
