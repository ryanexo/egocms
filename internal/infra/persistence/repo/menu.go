package repo

import (
    "context"
    "fmt"
    
    "dpcms/internal/app/dto"
    "dpcms/internal/infra/persistence/dbscope"
    "dpcms/internal/infra/persistence/model"
    "dpcms/internal/infra/persistence/query"
    
    "gorm.io/gen"
)

type MenuRepo struct {
    query *query.Query
}

func NewMenuRepo(persist *query.Query) *MenuRepo {
    return &MenuRepo{persist}
}

func (s *MenuRepo) Create(ctx context.Context, menu *model.Menu) error {
    return s.query.Menu.WithContext(ctx).Create(menu)
}

func (s *MenuRepo) CreateSubtree(ctx context.Context, id uint64, parentId uint64) error {
    return s.query.MenuContext.WithContext(ctx).CreateSubtree(id, parentId)
}

func (s *MenuRepo) Update(ctx context.Context, data *model.Menu) (gen.ResultInfo, error) {
    return s.query.Menu.WithContext(ctx).Where(s.query.Menu.ID.Eq(data.ID.Raw())).Updates(data)
}

func (s *MenuRepo) Move(ctx context.Context, fromNode uint64, toNode uint64) error {
    ctxDao := s.query.MenuContext
    catDao := s.query.Menu
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

func (s *MenuRepo) Delete(ctx context.Context, id uint64) error {
    ctxDao := s.query.MenuContext
    catDao := s.query.Menu
    sqlStr := `DELETE FROM %[1]s WHERE id IN ( SELECT d_id FROM ( SELECT t.%[4]s AS d_id FROM %[2]s AS t WHERE %[3]s = ? ) )`
    deleteCategorySql := fmt.Sprintf(
        sqlStr,
        catDao.TableName(),
        ctxDao.TableName(),
        ctxDao.Ancestor.ColumnName(),
        ctxDao.Descendant.ColumnName(),
    )
    
    err := s.query.Menu.WithContext(ctx).UnderlyingDB().Exec(deleteCategorySql, id).Error
    if err != nil {
        return err
    }
    
    err = s.query.MenuContext.WithContext(ctx).DropSubtree(id)
    if err != nil {
        return err
    }
    
    return nil
}

func (s *MenuRepo) FindByID(ctx context.Context, id uint64) (*model.Menu, error) {
    return s.query.Menu.WithContext(ctx).Where(s.query.Menu.ID.Eq(id)).First()
}

func (s *MenuRepo) FindByIDWithAncestor(ctx context.Context, ancestor uint64, descendant uint64) (*model.MenuContext, error) {
    ctxDao := s.query.MenuContext
    return ctxDao.WithContext(ctx).Where(ctxDao.Ancestor.Eq(ancestor), ctxDao.Descendant.Eq(descendant)).First()
}

func (s *MenuRepo) List(ctx context.Context, params dto.MenuListQueryParams) ([]*model.Menu, int64, error) {
    menuDAO := s.query.Menu
    q := menuDAO.WithContext(ctx)
    if params.ParentID != nil {
        q = q.Where(menuDAO.ParentID.Eq(params.ParentID.Raw()))
    }
    if params.Name != nil {
        q = q.Where(menuDAO.Name.Like("%" + *params.Name + "%"))
    }
    count, err := q.Count()
    if err != nil {
        return nil, 0, err
    }
    menuList, err := q.Scopes(dbscope.Paginate(params.PageNo, params.PageSize)).Find()
    if err != nil {
        return nil, 0, err
    }
    return menuList, count, nil
}
