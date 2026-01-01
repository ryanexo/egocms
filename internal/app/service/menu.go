package service

import (
    `context`
    `errors`
    `fmt`
    
    `dpcms/internal/app/dto`
    `dpcms/internal/app/erroz`
    `dpcms/internal/app/service/internal/common`
    `dpcms/internal/app/util/gormutil`
    `dpcms/internal/infra`
    `dpcms/internal/infra/persistence/datatype`
    `dpcms/internal/infra/persistence/dbscope`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/infra/persistence/query`
    
    `github.com/jinzhu/copier`
    `gorm.io/gorm`
)

type Menu struct {
    persist *query.Query
    db      *gorm.DB
}

func NewMenuService(i *infra.Infra) *Menu {
    return &Menu{persist: i.Query, db: i.DB}
}

func (srv Menu) Create(ctx context.Context, createParams dto.MenuCreateParams) (*model.Menu, error) {
    menu := model.Menu{}
    if err := copier.Copy(&menu, &createParams); err != nil {
        return nil, err
    }
    err := srv.persist.Transaction(func(tx *query.Query) error {
        queryCtx := tx.WithContext(ctx)
        err := queryCtx.Menu.Create(&menu)
        if err != nil {
            return err
        }
        
        err = queryCtx.MenuContext.CreateSubtree(menu.ID.Raw(), menu.ParentID.Raw())
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

func (srv Menu) Update(ctx context.Context, params dto.MenuUpdateParams) error {
    return srv.persist.Transaction(func(tx *query.Query) error {
        queryCtx := tx.WithContext(ctx)
        _, err := queryCtx.Menu.Where(tx.Menu.ID.Eq(params.ID.Raw())).First()
        if err != nil {
            return err
        }
        _, err = tx.WithContext(ctx).Category.Where(tx.Menu.ID.Eq(params.ID.Raw())).Updates(params)
        if err != nil {
            return err
        }
        return nil
    })
}

func (srv Menu) Delete(ctx context.Context, id datatype.SafeUint64) error {
    return srv.persist.Transaction(func(tx *query.Query) error {
        ctxDao := tx.MenuContext
        
        sqlStr := "DELETE FROM %[1]s WHERE id IN ( SELECT d_id FROM ( SELECT t.%[3]s AS d_id FROM %[1]s AS t WHERE %[2]s = ? ) )"
        deleteSql := fmt.Sprintf(sqlStr, tx.Menu.TableName(), ctxDao.Ancestor.ColumnName(), ctxDao.Descendant.ColumnName())
        
        err := srv.db.WithContext(ctx).Exec(deleteSql, id).Error
        if err != nil {
            return err
        }
        err = tx.WithContext(ctx).CategoryContext.DropSubtree(id.Raw())
        if err != nil {
            return err
        }
        return nil
    })
}

func (srv Menu) Move(ctx context.Context, id datatype.SafeUint64, target datatype.SafeUint64) error {
    return srv.persist.Transaction(func(tx *query.Query) error {
        queryCtx := tx.WithContext(ctx)
        ctxDao := tx.MenuContext
        
        _, err := queryCtx.MenuContext.Where(ctxDao.Ancestor.Eq(id.Raw()), ctxDao.Descendant.Eq(target.Raw())).First()
        if err == nil {
            return erroz.MenuCircular.ToError()
        } else if !errors.Is(err, gorm.ErrRecordNotFound) {
            return err
        }
        
        err = queryCtx.MenuContext.UnbindRelationships(id.Raw())
        if err != nil {
            return err
        }
        err = queryCtx.MenuContext.ReBindRelationships(id.Raw(), target.Raw())
        if err != nil {
            return err
        }
        _, err = queryCtx.Menu.Where(tx.Menu.ID.Eq(id.Raw())).Update(tx.Menu.ParentID, target)
        if err != nil {
            return err
        }
        
        return nil
    })
}

func (srv Menu) FindByID(ctx context.Context, id datatype.SafeUint64) (*model.Menu, error) {
    q := srv.persist.Menu
    result, err := q.WithContext(ctx).Where(q.ID.Eq(id.Raw())).First()
    if err != nil {
        return nil, gormutil.ReplaceNotFoundError(err)
    }
    return result, nil
}

func (srv Menu) List(ctx context.Context, condition dto.MenuListQueryParams) (result *common.PaginatedResult[*model.Menu], err error) {
    menuDAO := srv.persist.Menu
    q := menuDAO.WithContext(ctx)
    if condition.Ancestor != nil {
        q = q.Where(menuDAO.ParentID.Eq(condition.Ancestor.Raw()))
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
