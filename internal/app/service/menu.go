package service

import (
    `context`
    `fmt`
    
    `dpcms/internal/app/helper/dbscope`
    `dpcms/internal/app/helper/gormhelper`
    `dpcms/internal/app/service/internal/common`
    `dpcms/internal/app/service/srvparams`
    `dpcms/internal/database/model`
    `dpcms/internal/database/query`
    `dpcms/internal/infra`
    
    `github.com/jinzhu/copier`
)

type MenuService struct {
    query *query.Query
    infra *infra.Infra
}

func NewMenuService(infra *infra.Infra) *MenuService {
    return &MenuService{query: query.Use(infra.DB), infra: infra}
}

func (srv MenuService) Create(ctx context.Context, createParams srvparams.MenuCreateParams) (*model.Menu, error) {
    menu := model.Menu{}
    if err := copier.Copy(&menu, &createParams); err != nil {
        return nil, err
    }
    err := srv.query.Transaction(func(tx *query.Query) error {
        queryCtx := tx.WithContext(ctx)
        err := queryCtx.Menu.Create(&menu)
        if err != nil {
            return err
        }
        
        err = queryCtx.MenuContext.CreateSubtree(menu.ID, menu.ParentID)
        if err != nil {
            return err
        }
        return nil
    })
    if err != nil {
        return nil, err
    }
    return &menu, err
}

func (srv MenuService) Update(ctx context.Context, params srvparams.MenuUpdateParams) error {
    return srv.query.Transaction(func(tx *query.Query) error {
        queryCtx := tx.WithContext(ctx)
        _, err := queryCtx.Menu.Where(srv.query.Menu.ID.Eq(params.ID)).First()
        if err != nil {
            return err
        }
        _, err = srv.query.WithContext(ctx).Category.Where(srv.query.Menu.ID.Eq(params.ID)).Updates(params)
        if err != nil {
            return err
        }
        return nil
    })
}

func (srv MenuService) Delete(ctx context.Context, id uint64) error {
    return srv.query.Transaction(func(tx *query.Query) error {
        ctxDao := srv.query.MenuContext
        
        sqlStr := "DELETE FROM %[1]s WHERE id IN ( SELECT d_id FROM ( SELECT t.%[3]s AS d_id FROM %[1]s AS t WHERE %[2]s = ? ) )"
        deleteSql := fmt.Sprintf(sqlStr, srv.query.Menu.TableName(), ctxDao.Ancestor.ColumnName(), ctxDao.Descendant.ColumnName())
        
        err := srv.infra.DB.WithContext(ctx).Exec(deleteSql, id).Error
        if err != nil {
            return err
        }
        err = tx.WithContext(ctx).CategoryContext.DropSubtree(id)
        if err != nil {
            return err
        }
        return nil
    })
}

func (srv MenuService) Move(ctx context.Context, id uint64, target uint64) error {
    return srv.query.Transaction(func(tx *query.Query) error {
        queryCtx := tx.WithContext(ctx)
        
        err := queryCtx.MenuContext.UnbindRelationships(id)
        if err != nil {
            return err
        }
        err = queryCtx.MenuContext.ReBindRelationships(id, target)
        if err != nil {
            return err
        }
        _, err = queryCtx.Menu.Where(srv.query.Menu.ID.Eq(id)).Update(srv.query.Menu.ParentID, target)
        if err != nil {
            return err
        }
        
        return nil
    })
}

func (srv MenuService) FindByID(ctx context.Context, id uint64) (*model.Menu, error) {
    q := srv.query.Menu
    result, err := q.WithContext(ctx).Where(q.ID.Eq(id)).First()
    if err != nil {
        return nil, gormhelper.ReplaceNotFoundError(err)
    }
    return result, nil
}

func (srv MenuService) List(ctx context.Context, condition srvparams.MenuListQueryParams) (result *common.PaginatedResult[*model.Menu], err error) {
    menuDAO := srv.query.Menu
    q := menuDAO.WithContext(ctx)
    if condition.Ancestor != nil {
        q = q.Where(menuDAO.ParentID.Eq(*condition.Ancestor))
    }
    if condition.Name != nil {
        q = q.Where(menuDAO.Name.Like(*condition.Name))
    }
    count, err := q.Count()
    if err != nil {
        return
    }
    menuList, err := q.Scopes(dbscope.Paginate(condition.PageNo, condition.PageSize)).Find()
    if err != nil {
        return
    }
    result = &common.PaginatedResult[*model.Menu]{
        Total:    count,
        List:     menuList,
        PageNo:   condition.PageNo,
        PageSize: condition.PageSize,
    }
    return
}
