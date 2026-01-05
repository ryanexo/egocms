package repo

import (
    `context`
    `fmt`
    
    `dpcms/internal/app/dto`
    `dpcms/internal/infra/persistence/dbscope`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/infra/persistence/query`
    
    `gorm.io/gen`
)

type CategoryRepo struct {
    persist *query.Query
}

func NewCategoryRepo(persist *query.Query) CategoryRepo {
    return CategoryRepo{persist}
}

func (s CategoryRepo) Create(ctx context.Context, category *model.Category) error {
    return s.persist.Category.WithContext(ctx).Create(category)
}

func (s CategoryRepo) CreateSubtree(ctx context.Context, id uint64, parentId uint64) error {
    return s.persist.CategoryContext.WithContext(ctx).CreateSubtree(id, parentId)
}

func (s CategoryRepo) Update(ctx context.Context, category *model.Category) (gen.ResultInfo, error) {
    return s.persist.Category.WithContext(ctx).Updates(category)
}

func (s CategoryRepo) Move(ctx context.Context, fromNode uint64, toNode uint64) error {
    ctxDao := s.persist.CategoryContext
    catDao := s.persist.Category
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

func (s CategoryRepo) Delete(ctx context.Context, id uint64) error {
    ctxDao := s.persist.CategoryContext
    catDao := s.persist.Category
    seoDao := s.persist.CategorySeo
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
    
    err := s.persist.Category.WithContext(ctx).UnderlyingDB().Exec(deleteCategorySql, id).Error
    if err != nil {
        return err
    }
    
    err = s.persist.CategorySeo.WithContext(ctx).UnderlyingDB().Exec(deleteSeoSql, id).Error
    if err != nil {
        return err
    }
    
    err = s.persist.CategoryContext.WithContext(ctx).DropSubtree(id)
    if err != nil {
        return err
    }
    
    return nil
}

func (s CategoryRepo) FindByID(ctx context.Context, id uint64) (*model.Category, error) {
    return s.persist.Category.WithContext(ctx).Preload(s.persist.Category.SEO).Where(s.persist.Category.ID.Eq(id)).First()
}

func (s CategoryRepo) FindByIDWithAncestor(ctx context.Context, ancestor uint64, descendant uint64) (*model.CategoryContext, error) {
    ctxDao := s.persist.CategoryContext
    return ctxDao.WithContext(ctx).Where(ctxDao.Ancestor.Eq(ancestor), ctxDao.Descendant.Eq(descendant)).First()
}

func (s CategoryRepo) List(ctx context.Context, params dto.CategoryListParams) ([]*model.Category, int64, error) {
    dao := s.persist.Category
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
