package service

import (
    `context`
    `errors`
    
    `dpcms/internal/database/model`
    `dpcms/internal/database/query`
    `dpcms/internal/http/errors/menu_error`
    `dpcms/internal/infra`
    `dpcms/internal/utils/dbscopes`
    `gorm.io/gen`
    `gorm.io/gorm`
    `gorm.io/gorm/clause`
)

type MenuService struct {
    query *query.Query
}

func NewMenuService(infra *infra.Infra) *MenuService {
    return &MenuService{query: query.Use(infra.DB)}
}

// Create
//
// ParentID为0时，视为根节点
func (srv *MenuService) Create(ctx context.Context, menu *model.Menu) error {
    return srv.query.Transaction(func(tx *query.Query) error {
        queryCtx := tx.WithContext(ctx)
        err := queryCtx.Menu.Create(menu)
        if err != nil {
            return err
        }
        
        err = queryCtx.MenuContext.CreateBranch(menu.ID, menu.ParentID)
        if err != nil {
            return err
        }
        
        return nil
    })
}

func (srv *MenuService) Update(ctx context.Context, menu *model.Menu) error {
    return srv.query.Transaction(func(tx *query.Query) error {
        queryCtx := tx.WithContext(ctx)
        ct, err := queryCtx.Menu.Where(srv.query.Menu.ID.Eq(menu.ID)).First()
        if err != nil {
            return err
        }
        _, err = srv.query.WithContext(ctx).Category.Updates(menu)
        if err != nil {
            return err
        }
        if ct.ParentID != menu.ParentID {
            return srv.Move(ctx, menu.ID, menu.ParentID)
        }
        return nil
    })
}

func (srv *MenuService) Delete(ctx context.Context, id int64) error {
    return srv.query.Transaction(func(tx *query.Query) error {
        queryCtx := tx.WithContext(ctx)
        
        menuDAO := tx.Menu
        
        _, err := queryCtx.Menu.Where(menuDAO.ID.Eq(id)).Delete()
        if err != nil {
            return err
        }
        
        err = queryCtx.MenuContext.RemoveBranch(id)
        if err != nil {
            return err
        }
        
        return nil
    })
}

func (srv *MenuService) Move(ctx context.Context, id int64, parent int64) error {
    return srv.query.Transaction(func(tx *query.Query) error {
        queryCtx := tx.WithContext(ctx)
        
        menuDAO := tx.Menu
        menuCtxDAO := tx.MenuContext
        
        _, err := queryCtx.MenuContext.
            Where(menuCtxDAO.Ancestor.Eq(id)).
            Where(menuCtxDAO.Descendant.Eq(parent)).
            First()
        
        if err == nil {
            currentCategory, err := queryCtx.Menu.Where(menuCtxDAO.ID.Eq(id)).First()
            if err != nil {
                return err
            }
            targetCategory, err := queryCtx.Menu.Where(menuDAO.ID.Eq(parent)).First()
            if err != nil {
                return err
            }
            return menu_error.ErrCircularReferenceWhenMove.Format(targetCategory.Name, currentCategory.Name).ToError()
        } else if !errors.Is(err, gorm.ErrRecordNotFound) {
            return err
        }
        
        _, err = queryCtx.Menu.Where(menuDAO.ID.Eq(id)).Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).Find()
        
        _, err = queryCtx.Menu.Where(menuDAO.ID.Eq(id)).Update(menuDAO.ParentID, parent)
        if err != nil {
            return err
        }
        
        descendants, err := queryCtx.MenuContext.FindDescendantByAncestor(id)
        if err != nil {
            return err
        }
        
        err = queryCtx.MenuContext.RemoveBranch(id)
        if err != nil {
            return err
        }
        
        err = queryCtx.MenuContext.CreateBranch(id, parent)
        if err != nil {
            return err
        }
        
        descendantMap := make(map[int64]*model.MenuContext)
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
            return queryCtx.MenuContext.CreateBranch(s.Descendant, s.Parent)
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

func (srv *MenuService) FindByID(ctx context.Context, id int64) (*model.Menu, error) {
    q := srv.query.Menu
    return q.WithContext(ctx).Where(q.ID.Eq(id)).First()
}

func (srv *MenuService) FindDescendant(ctx context.Context, id int64, distance int64) ([]*model.Menu, error) {
    queryCtx := srv.query.WithContext(ctx)
    menuCtxDAO := srv.query.MenuContext
    cond := menuCtxDAO.Distance.Lt(distance)
    if distance == 0 {
        cond = menuCtxDAO.Distance.Gt(0)
    }
    descendants := queryCtx.MenuContext.Select(menuCtxDAO.ID).Where(cond).Where(menuCtxDAO.Ancestor.Eq(id))
    return queryCtx.Menu.Where(gen.Exists(descendants)).Find()
}

func (srv *MenuService) ListRootNodes(ctx context.Context, pageNo int, pageSize int) ([]*model.Menu, error) {
    return srv.query.Menu.WithContext(ctx).Scopes(dbscopes.Paginate(pageNo, pageSize)).Where(srv.query.Menu.ParentID.Eq(0)).Find()
}

func (srv *MenuService) ListNodesByParentID(ctx context.Context, id int64, pageSize int, pageNo int) ([]*model.Menu, error) {
    q := srv.query.Menu
    return q.WithContext(ctx).Where(q.ParentID.Eq(id)).Scopes(dbscopes.Paginate(pageNo, pageSize)).Find()
}
