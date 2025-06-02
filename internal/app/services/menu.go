package services

import (
    `context`
    `errors`
    
    `dpcms/internal/app/errors/menu_error`
    `dpcms/internal/app/services/types`
    `dpcms/internal/app/services/types/menu`
    `dpcms/internal/database/model`
    `dpcms/internal/database/query`
    `dpcms/internal/infra`
    `dpcms/internal/utils/dbscopes`
    `github.com/jinzhu/copier`
    `gorm.io/gen`
    `gorm.io/gorm`
    `gorm.io/gorm/clause`
)

type MenuService struct {
    query *query.Query
}

func NewMenuService(infra *infra.Infra) MenuService {
    return MenuService{query: query.Use(infra.DB)}
}

// Create
//
// ParentID为0时，视为根节点
func (srv MenuService) Create(ctx context.Context, params menu.CreateParams) (result menu.Detail, err error) {
    m := model.Menu{Visible: true}
    err = copier.Copy(&m, &params)
    if err != nil {
        return
    }
    err = srv.query.Transaction(func(tx *query.Query) error {
        queryCtx := tx.WithContext(ctx)
        err := queryCtx.Menu.Create(&m)
        if err != nil {
            return err
        }
        
        err = queryCtx.MenuContext.CreateBranch(m.ID, m.ParentID)
        if err != nil {
            return err
        }
        return nil
    })
    err = copier.Copy(&result, &m)
    return
}

func (srv MenuService) Update(ctx context.Context, params menu.UpdateParams) (result menu.Detail, err error) {
    menuModel := model.Menu{}
    err = copier.Copy(&menuModel, &params)
    if err != nil {
        return
    }
    err = srv.query.Transaction(func(tx *query.Query) error {
        queryCtx := tx.WithContext(ctx)
        m, dbErr := queryCtx.Menu.Where(srv.query.Menu.ID.Eq(params.ID)).First()
        if dbErr != nil {
            return dbErr
        }
        _, dbErr = srv.query.WithContext(ctx).Category.Where(srv.query.Menu.ID.Eq(params.ID)).Updates(menuModel)
        if dbErr != nil {
            return dbErr
        }
        if m.ParentID != params.ParentID {
            if dbErr = srv.Move(ctx, m.ID, params.ParentID); dbErr != nil {
                return dbErr
            }
        }
        menuModel = *m
        return nil
    })
    if err != nil {
        return
    }
    err = copier.Copy(&result, &menuModel)
    return
}

func (srv MenuService) Delete(ctx context.Context, id int64) error {
    return srv.query.Transaction(func(tx *query.Query) error {
        queryCtx := tx.WithContext(ctx)
        
        menuDAO := srv.query.Menu
        
        _, err := queryCtx.Menu.Where(menuDAO.ID.Eq(id)).Delete()
        descendants, err := srv.FindDescendant(ctx, id, 0)
        if err != nil {
            return err
        }
        
        idList := make([]int64, 0, len(descendants))
        for _, descendant := range descendants {
            idList = append(idList, descendant.ID)
        }
        for i := 0; i < len(idList); i += 100 {
            items := idList[i : i+100]
            _, err = queryCtx.Menu.Where(srv.query.Menu.ID.In(items...)).Delete()
            if err != nil {
                return err
            }
        }
        
        err = queryCtx.MenuContext.RemoveBranch(id)
        if err != nil {
            return err
        }
        
        return nil
    })
}

func (srv MenuService) Move(ctx context.Context, id int64, parent int64) error {
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

func (srv MenuService) FindByID(ctx context.Context, id int64) (*model.Menu, error) {
    q := srv.query.Menu
    return q.WithContext(ctx).Where(q.ID.Eq(id)).First()
}

func (srv MenuService) FindDescendant(ctx context.Context, id int64, distance int64) ([]*model.Menu, error) {
    queryCtx := srv.query.WithContext(ctx)
    menuCtxDAO := srv.query.MenuContext
    cond := menuCtxDAO.Distance.Lt(distance)
    if distance == 0 {
        cond = menuCtxDAO.Distance.Gt(0)
    }
    descendants := queryCtx.MenuContext.Select(menuCtxDAO.ID).Where(cond).Where(menuCtxDAO.Ancestor.Eq(id))
    return queryCtx.Menu.Where(gen.Exists(descendants)).Find()
}

func (srv MenuService) List(ctx context.Context, condition *menu.RetrieveListParams) (result types.Pagination[menu.Detail], err error) {
    menuDAO := srv.query.Menu
    q := menuDAO.WithContext(ctx).Debug()
    if condition.Ancestor != nil {
        if condition.Recursive {
            ctxDAO := srv.query.MenuContext
            q = q.Where(menuDAO.Columns(menuDAO.ID).In(
                ctxDAO.WithContext(ctx).Select(ctxDAO.Descendant).Where(ctxDAO.Ancestor.Eq(*condition.Ancestor)),
            ))
        } else {
            q = q.Where(menuDAO.ParentID.Eq(*condition.Ancestor))
        }
    }
    if condition.Name != nil {
        q = q.Where(menuDAO.Name.Like(*condition.Name))
    }
    count, err := q.Count()
    if err != nil {
        return
    }
    menuList, err := q.Scopes(dbscopes.Paginate(condition.PageNo, condition.PageSize)).Find()
    if err != nil {
        return
    }
    details := make([]menu.Detail, len(menuList), len(menuList))
    err = copier.Copy(&details, menuList)
    if err != nil {
        return
    }
    result = types.Pagination[menu.Detail]{
        Total:    count,
        List:     details,
        PageNo:   condition.PageNo,
        PageSize: condition.PageSize,
    }
    return
}
