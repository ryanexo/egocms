package service

import (
    "context"
    `errors`
    
    `dpcms/internal/app/erroz`
    `dpcms/internal/app/helper/dbscope`
    `dpcms/internal/app/service/types/category`
    `dpcms/internal/database/model`
    `dpcms/internal/database/query`
    `dpcms/internal/infra`
    
    `gorm.io/gorm`
    `gorm.io/gorm/clause`
)

type CategoryService struct {
    query *query.Query
}

func NewCategoryCategory(infra *infra.Infra) CategoryService {
    return CategoryService{query: query.Use(infra.DB)}
}

// Create
//
// ParentID为0时，视为根节点
func (srv CategoryService) Create(ctx context.Context, category category.Detail) (err error) {
    return srv.query.Transaction(func(tx *query.Query) error {
        queryCtx := tx.WithContext(ctx)
        err := queryCtx.Category.Create(&model.Category{
            ParentID: category.ParentID,
            Sequence: category.Sequence,
            Name:     category.Name,
            Path:     category.Path,
            Type:     category.Type,
            Display:  category.Display,
            SEO:      category.SEO,
            Children: nil,
        })
        if err != nil {
            return err
        }
        
        err = queryCtx.CategoryContext.CreateBranch(category.ID, category.ParentID)
        if err != nil {
            return err
        }
        
        err = srv.query.Category.SEO.Model(category).Append(category.SEO)
        if err != nil {
            return err
        }
        
        return nil
    })
}

func (srv CategoryService) Update(ctx context.Context, category *model.Category) error {
    return srv.query.Transaction(func(tx *query.Query) error {
        queryCtx := tx.WithContext(ctx)
        ct, err := queryCtx.Category.Where(srv.query.Category.ID.Eq(category.ID)).First()
        if err != nil {
            return err
        }
        _, err = srv.query.WithContext(ctx).Category.Updates(category)
        if err != nil {
            return err
        }
        if ct.ParentID != category.ParentID {
            return srv.Move(ctx, category.ID, category.ParentID)
        }
        return nil
    })
}

func (srv CategoryService) Delete(ctx context.Context, id int64) error {
    return srv.query.Transaction(func(tx *query.Query) error {
        queryCtx := tx.WithContext(ctx)
        
        categoryDAO := tx.Category
        
        _, err := queryCtx.Category.Where(categoryDAO.ID.Eq(id)).Delete()
        if err != nil {
            return err
        }
        
        err = queryCtx.CategoryContext.RemoveBranch(id)
        if err != nil {
            return err
        }
        
        return nil
    })
}

func (srv CategoryService) Move(ctx context.Context, id int64, parent int64) error {
    return srv.query.Transaction(func(tx *query.Query) error {
        queryCtx := tx.WithContext(ctx)
        
        categoryDAO := tx.Category
        categoryCtxDAO := tx.CategoryContext
        
        _, err := queryCtx.CategoryContext.
            Where(categoryCtxDAO.Ancestor.Eq(id)).
            Where(categoryCtxDAO.Descendant.Eq(parent)).
            First()
        
        if err == nil {
            currentCategory, err := queryCtx.Category.Where(categoryCtxDAO.ID.Eq(id)).First()
            if err != nil {
                return err
            }
            targetCategory, err := queryCtx.Category.Where(categoryDAO.ID.Eq(parent)).First()
            if err != nil {
                return err
            }
            return erroz.ErrCategoryCircularReferenceWhenMove.Format(targetCategory.Name, currentCategory.Name).ToError()
        } else if !errors.Is(err, gorm.ErrRecordNotFound) {
            return err
        }
        
        _, err = queryCtx.Category.Where(categoryDAO.ID.Eq(id)).Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).Find()
        
        _, err = queryCtx.Category.Where(categoryDAO.ID.Eq(id)).Update(categoryDAO.ParentID, parent)
        if err != nil {
            return err
        }
        
        descendants, err := queryCtx.CategoryContext.FindDescendantByAncestor(id)
        if err != nil {
            return err
        }
        
        err = queryCtx.CategoryContext.RemoveBranch(id)
        if err != nil {
            return err
        }
        
        err = queryCtx.CategoryContext.CreateBranch(id, parent)
        if err != nil {
            return err
        }
        
        descendantMap := make(map[int64]*model.CategoryContext)
        processed := make(map[int64]bool)
        
        var createBranch func(id int64) error
        
        createBranch = func(id int64) error {
            if processed[id] {
                return nil
            }
            s, ok := descendantMap[id]
            if !ok {
                return nil
            }
            if s.Parent != id && !processed[s.Parent] {
                if err := createBranch(s.Parent); err != nil {
                    return err
                }
            }
            return queryCtx.CategoryContext.CreateBranch(s.Descendant, s.Parent)
        }
        
        for _, descendant := range descendants {
            descendantMap[descendant.Descendant] = descendant
        }
        
        for _, s := range descendants {
            err = createBranch(s.Descendant)
            if err != nil {
                return err
            }
        }
        
        return nil
    })
}

func (srv CategoryService) FindByID(ctx context.Context, id int64) (*model.Category, error) {
    q := srv.query.Category
    return q.WithContext(ctx).Where(q.ID.Eq(id)).First()
}

func (srv CategoryService) ListRootNodes(ctx context.Context, pageNo int, pageSize int) ([]*model.Category, error) {
    return srv.query.Category.WithContext(ctx).Scopes(dbscope.Paginate(pageNo, pageSize)).Where(srv.query.Category.ParentID.Eq(0)).Find()
}

func (srv CategoryService) ListNodesByParentID(ctx context.Context, id int64, pageSize int, pageNo int) ([]*model.Category, error) {
    q := srv.query.Category
    return q.WithContext(ctx).Where(q.ParentID.Eq(id)).Scopes(dbscope.Paginate(pageNo, pageSize)).Find()
}
