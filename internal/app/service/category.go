package service

import (
    "context"
    `errors`
    `fmt`
    
    `dpcms/internal/app/dto`
    `dpcms/internal/app/erroz`
    `dpcms/internal/app/helper/gormhelper`
    `dpcms/internal/app/service/internal/common`
    `dpcms/internal/infra`
    `dpcms/internal/infra/persistence/dbscope`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/infra/persistence/query`
    
    `github.com/jinzhu/copier`
    `gorm.io/gorm`
)

type Category struct {
    persist *query.Query
    db      *gorm.DB
}

func NewCategoryCategory(i *infra.Infra) *Category {
    return &Category{persist: i.Query, db: i.DB}
}

func (srv Category) Create(ctx context.Context, createParams dto.CategoryCreateParams) (*model.Category, error) {
    category := model.Category{}
    if err := copier.Copy(&category, &createParams); err != nil {
        return nil, err
    }
    err := srv.persist.Transaction(func(tx *query.Query) error {
        queryCtx := tx.WithContext(ctx)
        txErr := queryCtx.Category.Create(&category)
        if txErr != nil {
            return txErr
        }
        
        txErr = queryCtx.CategoryContext.CreateSubtree(category.ID, category.ParentID)
        if txErr != nil {
            return txErr
        }
        
        return nil
    })
    if err != nil {
        return nil, err
    }
    return &category, nil
}

func (srv Category) Update(ctx context.Context, category dto.CategoryUpdateParams) error {
    return srv.persist.Transaction(func(tx *query.Query) error {
        queryCtx := tx.WithContext(ctx)
        _, err := queryCtx.Category.Where(srv.persist.Category.ID.Eq(category.ID)).First()
        if err != nil {
            return err
        }
        _, err = srv.persist.WithContext(ctx).Category.Updates(category)
        return err
    })
}

func (srv Category) Delete(ctx context.Context, id uint64) error {
    return srv.persist.Transaction(func(tx *query.Query) error {
        ctxDao := srv.persist.CategoryContext
        
        sqlStr := "DELETE FROM %[1]s WHERE id IN ( SELECT d_id FROM ( SELECT t.%[3]s AS d_id FROM %[1]s AS t WHERE %[2]s = ? ) )"
        deleteSql := fmt.Sprintf(sqlStr, srv.persist.Category.TableName(), ctxDao.Ancestor.ColumnName(), ctxDao.Descendant.ColumnName())
        
        err := srv.db.WithContext(ctx).Exec(deleteSql, id).Error
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

func (srv Category) Move(ctx context.Context, id uint64, target uint64) error {
    return srv.persist.Transaction(func(tx *query.Query) error {
        queryCtx := tx.WithContext(ctx)
        ctxDao := srv.persist.CategoryContext
        
        _, err := queryCtx.CategoryContext.Where(ctxDao.Ancestor.Eq(id), ctxDao.Descendant.Eq(target)).First()
        if err == nil {
            return erroz.CategoryCircularReferenceWhenMove.ToError()
        } else if !errors.Is(err, gorm.ErrRecordNotFound) {
            return err
        }
        
        err = queryCtx.CategoryContext.UnbindRelationships(id)
        if err != nil {
            return err
        }
        err = queryCtx.CategoryContext.ReBindRelationships(id, target)
        if err != nil {
            return err
        }
        _, err = queryCtx.Category.Where(srv.persist.Category.ID.Eq(id)).Update(srv.persist.Category.ParentID, target)
        if err != nil {
            return err
        }
        
        return nil
    })
}

func (srv Category) FindByID(ctx context.Context, id uint64) (*model.Category, error) {
    q := srv.persist.Category
    result, err := q.WithContext(ctx).Where(q.ID.Eq(id)).First()
    if err != nil {
        return nil, gormhelper.ReplaceNotFoundError(err)
    }
    return result, nil
}

func (srv Category) ListRootNodes(ctx context.Context, pageNo int, pageSize int) ([]*model.Category, error) {
    return srv.persist.Category.WithContext(ctx).Scopes(dbscope.Paginate(pageNo, pageSize)).Where(srv.persist.Category.ParentID.Eq(0)).Find()
}

func (srv Category) ListNodesByParentID(ctx context.Context, id uint64, pageSize int, pageNo int) ([]*model.Category, error) {
    q := srv.persist.Category
    return q.WithContext(ctx).Where(q.ParentID.Eq(id)).Scopes(dbscope.Paginate(pageNo, pageSize)).Find()
}

func (srv Category) List(ctx context.Context, params dto.CategoryListParams) (*common.PaginatedResult[*model.Category], error) {
    dao := srv.persist.Category
    q := dao.WithContext(ctx).Scopes(dbscope.Paginate(params.PageNo, params.PageSize))
    if params.ID != nil {
        q = q.Where(dao.ID.Eq(*params.ID))
    }
    if params.ParentID != nil {
        q = q.Where(dao.ParentID.Eq(*params.ParentID))
    }
    if params.Type != nil {
        q = q.Where(dao.Type.Eq(*params.Type))
    }
    if params.Name != nil {
        q = q.Where(dao.Name.Eq(*params.Name))
    }
    if params.Path != nil {
        q = q.Where(dao.Path.Eq(*params.Path))
    }
    if params.Display != nil {
        q = q.Where(dao.Display.Eq(*params.Display))
    }
    
    total, err := q.Count()
    if err != nil {
        return nil, err
    }
    result, err := q.Find()
    if err != nil {
        return nil, err
    }
    
    return &common.PaginatedResult[*model.Category]{
        Total:    total,
        PageSize: params.PageSize,
        PageNo:   params.PageNo,
        List:     result,
    }, nil
}
