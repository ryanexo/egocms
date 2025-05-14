package repository

import (
    "context"
    "errors"
    
    `dpcms/model`
    `dpcms/model/query`
    "gorm.io/gen"
    "gorm.io/gorm"
    `gorm.io/gorm/clause`
)

var CircularReferenceError = errors.New("不能将节点移动至自身子节点或后代节点下，将引发循环引用")

type Category struct {
    query *query.Query
}

// Insert
//
// ancestor为0时，视为根节点
func (repo Category) Insert(ctx context.Context, category *model.Category) error {
    return repo.query.Transaction(func(tx *query.Query) error {
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

func (repo Category) Find(ctx context.Context, id uint) (*model.Category, error) {
    return repo.query.Category.WithContext(ctx).Where(repo.query.Category.ID.Eq(id)).First()
}

func (repo Category) FindAll(ctx context.Context) ([]*model.Category, error) {
    return repo.query.Category.WithContext(ctx).Find()
}

// Move
//
// 移动节点会重构被移动节点及其所有子节点的分支关系
//
// 父节点无法移动至自身的子节点
func (repo Category) Move(ctx context.Context, id uint, ancestor uint) error {
    return repo.query.Transaction(func(tx *query.Query) error {
        categoryDAO := tx.Category
        categoryCtxDAO := tx.CategoryContext
        
        _, err := categoryCtxDAO.WithContext(ctx).
            Where(categoryCtxDAO.Ancestor.Eq(id)).
            Where(categoryCtxDAO.Descendant.Eq(ancestor)).
            First()
        if err == nil {
            return CircularReferenceError
        } else if !errors.Is(err, gorm.ErrRecordNotFound) {
            return err
        }
        
        _, err = categoryDAO.WithContext(ctx).Where(categoryDAO.ID.Eq(id)).Update(categoryDAO.ParentID, ancestor)
        if err != nil {
            return err
        }
        
        descendants, err := repo.findDescendantRelation(ctx, tx, id)
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

// Delete
//
// 删除原理与Move同理
func (repo Category) Delete(ctx context.Context, id uint, reserveChildren bool) error {
    return repo.query.Transaction(func(tx *query.Query) error {
        categoryDAO := repo.query.Category
        categoryCtxDAO := repo.query.CategoryContext
        
        currentCategory, err := categoryDAO.WithContext(ctx).Where(categoryDAO.ID.Eq(id)).First()
        if err != nil {
            return err
        }
        _, err = categoryDAO.WithContext(ctx).Where(categoryDAO.ID.Eq(id)).Delete(currentCategory)
        if err != nil {
            return err
        }
        
        if reserveChildren {
            descendants, err := repo.findDescendantRelation(ctx, tx, id)
            if err != nil {
                return err
            }
            
            for _, descendant := range descendants {
                err = repo.Move(ctx, descendant.Descendant, currentCategory.ParentID)
                if err != nil {
                    return err
                }
            }
        }
        
        err = categoryCtxDAO.WithContext(ctx).RemoveBranch(id)
        if err != nil {
            return err
        }
        
        return nil
    })
}

func (repo Category) findDescendantRelation(ctx context.Context, tx *query.Query, parentID uint) ([]*model.CategoryContext, error) {
    categoryCtxDAO := tx.CategoryContext
    parentJoinQuery := categoryCtxDAO.As("b")
    // select ancestor, descendant, distance, b.ancestor as parent from ctx as a left join ctx as b on a.descendant = b.descendant and b.distance = 1 where a.ancestor = 1 and a.distance > 0
    return categoryCtxDAO.WithContext(ctx).Select(
        categoryCtxDAO.ALL,
        parentJoinQuery.Ancestor.As("parent"),
    ).LeftJoin(
        parentJoinQuery,
        categoryCtxDAO.Descendant.EqCol(parentJoinQuery.Descendant),
        parentJoinQuery.Distance.Eq(1),
    ).Where(
        categoryCtxDAO.Ancestor.Eq(parentID),
        categoryCtxDAO.Distance.Gt(0),
    ).Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).Find()
}

func (repo Category) FindDescendant(ctx context.Context, id uint, distance uint) ([]*model.Category, error) {
    categoryDAO := repo.query.Category
    categoryCtxDAO := repo.query.CategoryContext
    cond := categoryCtxDAO.Distance.Lt(distance)
    if distance == 0 {
        cond = categoryCtxDAO.Distance.Gt(0)
    }
    descendants := categoryCtxDAO.WithContext(ctx).Select(categoryCtxDAO.ID).Where(cond).Where(categoryCtxDAO.Ancestor.Eq(id))
    return categoryDAO.WithContext(ctx).Where(gen.Exists(descendants)).Find()
}

func (repo Category) FindAncestor(ctx context.Context, id uint, distance uint) ([]*model.Category, error) {
    categoryDAO := repo.query.Category
    categoryCtxDAO := repo.query.CategoryContext
    cond := categoryCtxDAO.Distance.Lt(distance)
    if distance == 0 {
        cond = categoryCtxDAO.Distance.Gt(0)
    }
    ancestors := categoryCtxDAO.WithContext(ctx).Select(categoryCtxDAO.ID).Where(cond).Where(categoryCtxDAO.Descendant.Eq(id))
    return categoryDAO.WithContext(ctx).Where(gen.Exists(ancestors)).Find()
}

func NewCategoryRepo(db *gorm.DB) *Category {
    return &Category{query: query.Use(db)}
}
