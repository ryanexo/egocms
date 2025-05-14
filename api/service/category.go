package service

import (
    "context"
    `errors`
    
    `dpcms/api/infra`
    `dpcms/erroz`
    `dpcms/model`
    `dpcms/model/query`
    `gorm.io/gorm`
    `gorm.io/gorm/clause`
)

type Category struct {
    query *query.Query
}

func NewCategory(infra infra.Infra) *Category {
    return &Category{query: query.Use(infra.DB)}
}

// Create
//
// ParentID为0时，视为根节点
func (srv *Category) Create(ctx context.Context, category *model.Category) error {
    return srv.query.Transaction(func(tx *query.Query) error {
        err := tx.Category.WithContext(ctx).Create(category)
        if err != nil {
            return err
        }
        
        err = tx.CategoryContext.WithContext(ctx).CreateBranch(category.ParentID, category.ID)
        if err != nil {
            return err
        }
        
        return nil
    })
}

func (srv *Category) Update(ctx context.Context, category *model.Category) error {
    _, err := srv.query.Category.WithContext(ctx).Updates(category)
    return err
}

func (srv *Category) Delete(ctx context.Context, id uint) error {
    return srv.query.Transaction(func(tx *query.Query) error {
        categoryDAO := tx.Category
        categoryCtxDAO := tx.CategoryContext
        
        _, err := categoryDAO.WithContext(ctx).Where(categoryDAO.ID.Eq(id)).First()
        if err != nil {
            return err
        }
        
        _, err = categoryDAO.WithContext(ctx).Where(categoryDAO.ID.Eq(id)).Delete()
        if err != nil {
            return err
        }
        
        descendants, err := srv.findDescendantByAncestorWithLock(ctx, tx, id)
        if err != nil {
            return err
        }
        
        for _, descendant := range descendants {
            err = categoryCtxDAO.WithContext(ctx).RemoveBranch(descendant.Descendant)
            if err != nil {
                return err
            }
        }
        
        err = categoryCtxDAO.WithContext(ctx).RemoveBranch(id)
        if err != nil {
            return err
        }
        
        return nil
    })
}

func (srv *Category) Move(ctx context.Context, id uint, ancestor uint) error {
    return srv.query.Transaction(func(tx *query.Query) error {
        categoryDAO := tx.Category
        categoryCtxDAO := tx.CategoryContext
        
        _, err := categoryCtxDAO.WithContext(ctx).
            Where(categoryCtxDAO.Ancestor.Eq(id)).
            Where(categoryCtxDAO.Descendant.Eq(ancestor)).
            First()
        if err == nil {
            return erroz.ErrCircularReferenceWhenMove.ToError()
        } else if !errors.Is(err, gorm.ErrRecordNotFound) {
            return err
        }
        
        _, err = categoryDAO.WithContext(ctx).Where(categoryDAO.ID.Eq(id)).Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).Find()
        
        _, err = categoryDAO.WithContext(ctx).Where(categoryDAO.ID.Eq(id)).Update(categoryDAO.ParentID, ancestor)
        if err != nil {
            return err
        }
        
        descendants, err := srv.findDescendantByAncestorWithLock(ctx, tx, id)
        if err != nil {
            return err
        }
        
        err = categoryCtxDAO.WithContext(ctx).RemoveBranch(id)
        if err != nil {
            return err
        }
        
        err = categoryCtxDAO.WithContext(ctx).CreateBranch(ancestor, id)
        if err != nil {
            return err
        }
        
        for _, s := range descendants {
            err = categoryCtxDAO.WithContext(ctx).RemoveBranch(s.Descendant)
            if err != nil {
                return err
            }
            err = categoryCtxDAO.WithContext(ctx).CreateBranch(s.Parent, s.Descendant)
            if err != nil {
                return err
            }
        }
        
        return nil
    })
}

func (srv *Category) FindByID(ctx context.Context, id uint) (*model.Category, error) {
    q := srv.query.Category
    return q.WithContext(ctx).Where(q.ID.Eq(id)).First()
}

func (srv *Category) FindByParentID(ctx context.Context, parentID uint) ([]*model.Category, error) {
    q := srv.query.Category
    return q.WithContext(ctx).Where(q.ParentID.Eq(parentID)).Find()
}

func (srv *Category) findDescendantByAncestorWithLock(ctx context.Context, q *query.Query, ancestor uint) ([]*model.CategoryContext, error) {
    ctxDAO := q.CategoryContext
    ctxChildrenDAO := q.CategoryContext.As("b")
    return ctxDAO.WithContext(ctx).Select(ctxDAO.ALL).LeftJoin(
        ctxChildrenDAO,
        ctxDAO.Descendant.EqCol(ctxChildrenDAO.Descendant),
        ctxChildrenDAO.Distance.Eq(1),
    ).Where(
        ctxDAO.Ancestor.Eq(ancestor),
        ctxDAO.Distance.Gt(0),
    ).Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).Find()
}
