package service

import (
    `context`
    `errors`
    
    `dpcms/internal/app/dto`
    `dpcms/internal/app/erroz`
    `dpcms/internal/infra`
    `dpcms/internal/infra/persistence/datatype`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/infra/persistence/query`
    `dpcms/internal/infra/persistence/repo`
    
    menuAssembler `dpcms/internal/app/assembler/menu`
    
    `gorm.io/gorm`
)

type Menu struct {
    persist *query.Query
}

func NewMenuService(i *infra.Infra) *Menu {
    return &Menu{persist: i.Query}
}

func (srv Menu) Create(ctx context.Context, params dto.MenuCreateParams) (*model.Menu, error) {
    menu := menuAssembler.BuildMenuCreateCommand(&params)
    err := srv.persist.Transaction(func(tx *query.Query) error {
        menuRepo := repo.NewMenuRepo(srv.persist)
        txErr := menuRepo.Create(ctx, menu)
        if txErr != nil {
            return txErr
        }
        txErr = menuRepo.CreateSubtree(ctx, menu.ID.Raw(), menu.ParentID.Raw())
        if txErr != nil {
            return txErr
        }
        return nil
    })
    if err != nil {
        return nil, err
    }
    return menu, nil
}

func (srv Menu) Update(ctx context.Context, params dto.MenuUpdateParams) error {
    menuRepo := repo.NewMenuRepo(srv.persist)
    _, err := menuRepo.FindByID(ctx, params.ID.Raw())
    if err != nil {
        return err
    }
    menu := menuAssembler.BuildMenuUpdateCommand(&params)
    _, err = menuRepo.Update(ctx, menu)
    return err
}

func (srv Menu) Delete(ctx context.Context, id datatype.SafeUint64) error {
    return srv.persist.Transaction(func(tx *query.Query) error {
        return repo.NewMenuRepo(srv.persist).Delete(ctx, id.Raw())
    })
}

func (srv Menu) Move(ctx context.Context, id datatype.SafeUint64, target datatype.SafeUint64) error {
    return srv.persist.Transaction(func(tx *query.Query) error {
        menuRepo := repo.NewMenuRepo(tx)
        _, err := menuRepo.FindByIDWithAncestor(ctx, id.Raw(), target.Raw())
        if err == nil {
            return erroz.CategoryCircular.ToError()
        } else if !errors.Is(err, gorm.ErrRecordNotFound) {
            return err
        }
        
        return menuRepo.Move(ctx, id.Raw(), target.Raw())
    })
}

func (srv Menu) FindByID(ctx context.Context, id datatype.SafeUint64) (*dto.Menu, error) {
    menu, err := repo.NewMenuRepo(srv.persist).FindByID(ctx, id.Raw())
    if err != nil {
        return nil, err
    }
    return menuAssembler.BuildMenuDTO(menu), nil
}

func (srv Menu) List(ctx context.Context, condition dto.MenuListQueryParams) (*dto.PaginatedResult[*dto.Menu], error) {
    data, total, err := repo.NewMenuRepo(srv.persist).List(ctx, condition)
    if err != nil {
        return nil, err
    }
    result := &dto.PaginatedResult[*dto.Menu]{
        Total:    total,
        List:     menuAssembler.BuildMenuListDTO(data),
        PageNo:   condition.PageNo,
        PageSize: condition.PageSize,
    }
    return result, nil
}
