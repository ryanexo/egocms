package provider

import (
    "context"
    "fmt"
    
    `cms/internal/domain/category/internal/dto`
    `cms/internal/infra/persist/contract`
    "cms/internal/infra/persist/model"
    "cms/internal/infra/persist/query"
    "cms/internal/infra/persist/scope"
)

type categoryRepo struct {
    query *query.Query
}

func NewCategoryRepo(persist *query.Query) contract.CategoryRepo {
    return &categoryRepo{persist}
}

func (s *categoryRepo) CloneWithQuery(q *query.Query) contract.CategoryRepo {
    return NewCategoryRepo(q)
}

func (s *categoryRepo) Create(ctx context.Context, category *model.Category) error {
    return s.query.Category.WithContext(ctx).Create(category)
}

func (s *categoryRepo) CreateSubtree(ctx context.Context, id uint64, parentID uint64) error {
    return s.query.CategoryContext.WithContext(ctx).CreateSubtree(id, parentID)
}

func (s *categoryRepo) Update(ctx context.Context, data *model.Category) (gen.ResultInfo, error) {
    return s.query.Category.WithContext(ctx).Where(s.query.Category.ID.Eq(data.ID.Raw())).Updates(data)
}

func (s *categoryRepo) Move(ctx context.Context, fromNode uint64, toNode uint64) error {
    ctxDao := s.query.CategoryContext
    catDao := s.query.Category
    err := ctxDao.WithContext(ctx).UnbindRelationships(fromNode)
    if err != nil {
        return err
    }
    err = ctxDao.WithContext(ctx).ReBindRelationships(fromNode, toNode)
    if err != nil {
        return err
    }
    _, err = catDao.WithContext(ctx).Where(catDao.ID.Eq(fromNode)).Update(catDao.ParentID, toNode)
    return err
}

func (s *categoryRepo) Delete(ctx context.Context, id uint64) error {
    ctxDao := s.query.CategoryContext
    catDao := s.query.Category
    seoDao := s.query.CategorySeo
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
    
    err := s.query.Category.WithContext(ctx).UnderlyingDB().Exec(deleteCategorySql, id).Error
    if err != nil {
        return err
    }
    
    err = s.query.CategorySeo.WithContext(ctx).UnderlyingDB().Exec(deleteSeoSql, id).Error
    if err != nil {
        return err
    }
    
    err = s.query.CategoryContext.WithContext(ctx).DropSubtree(id)
    if err != nil {
        return err
    }
    
    return nil
}

func (s *categoryRepo) FindByID(ctx context.Context, id uint64) (*model.Category, error) {
    return s.query.Category.WithContext(ctx).Preload(s.query.Category.SEO).Where(s.query.Category.ID.Eq(id)).First()
}

func (s *categoryRepo) FindByIDWithAncestor(ctx context.Context, ancestor uint64, descendant uint64) (*model.CategoryContext, error) {
    ctxDao := s.query.CategoryContext
    return ctxDao.WithContext(ctx).Where(ctxDao.Ancestor.Eq(ancestor), ctxDao.Descendant.Eq(descendant)).First()
}

func (s *categoryRepo) List(ctx context.Context, params dto.CategoryListParams) ([]*model.Category, int64, error) {
    dao := s.query.Category
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
