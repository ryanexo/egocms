package category

import (
    "context"
    "fmt"
    
    `cms/internal/app/category/internal`
    `cms/internal/infra/store/closuretable`
    `cms/internal/infra/store/gorm/dbscope`
    `cms/internal/infra/store/gorm/gquery`
    `cms/internal/modules/category/contract`
    `cms/internal/modules/category/internal/dto`
    
    "gorm.io/gen"
)

type categoryRepo struct {
    query   *gquery.Query
    closure closuretable.ClosureTable
}

func NewCategoryRepo(q *gquery.Query) (contract.CategoryRepo, error) {
    db, err := q.UnderlyingDB().DB()
    if err != nil {
        return nil, err
    }
    closure, err := closuretable.NewClosureTable(db, closuretable.WithModuleName("category"))
    if err != nil {
        return nil, err
    }
    return &categoryRepo{q, closure}, nil
}

func (s *categoryRepo) Create(ctx context.Context, category *internal.Category) error {
    return s.query.Category.WithContext(ctx).Create(category)
}

func (s *categoryRepo) CreateSubtree(ctx context.Context, id uint64, parentID uint64) error {
    return s.query.CategoryContext.WithContext(ctx).CreateSubtree(id, parentID)
}

func (s *categoryRepo) Update(ctx context.Context, data *internal.Category) (gen.ResultInfo, error) {
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

func (s *categoryRepo) FindByID(ctx context.Context, id uint64) (*internal.Category, error) {
    return s.query.Category.WithContext(ctx).Preload(s.query.Category.SEO).Where(s.query.Category.ID.Eq(id)).First()
}

func (s *categoryRepo) FindByIDWithAncestor(ctx context.Context, ancestor uint64, descendant uint64) (*internal.CategoryContext, error) {
    ctxDao := s.query.CategoryContext
    return ctxDao.WithContext(ctx).Where(ctxDao.Ancestor.Eq(ancestor), ctxDao.Descendant.Eq(descendant)).First()
}

func (s *categoryRepo) List(ctx context.Context, params dto.CategoryListParams) ([]*internal.Category, int64, error) {
    dao := s.query.Category
    q := dao.WithContext(ctx).Scopes(dbscope.Paginate(params.PageNo, params.PageSize))
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
