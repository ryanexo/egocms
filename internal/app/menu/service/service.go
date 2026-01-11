package service

import (
    `context`
    `errors`
    
    `dpcms/internal/app/menu/errno`
    `dpcms/internal/app/menu/internal/assembler`
    `dpcms/internal/app/menu/internal/dto`
    `dpcms/internal/infra/persistence/contract`
    `dpcms/internal/infra/persistence/datatype`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/infra/persistence/query`
    `dpcms/internal/types`
    
    `gorm.io/gorm`
)

type MenuService struct {
    txManager contract.TxManager
    repo      MenuRepo
}

func NewMenuService(txManager contract.TxManager, repo MenuRepo) *MenuService {
    return &MenuService{
        txManager: txManager,
        repo:      repo,
    }
}

func (s MenuService) Create(ctx context.Context, params dto.MenuCreateParams) (*model.Menu, error) {
    menu := assembler.BuildMenuCreateCommand(&params)
    err := s.txManager.Transaction(func(tx *query.Query) error {
        menuRepo := s.repo.CloneWithQuery(tx)
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

func (s MenuService) Update(ctx context.Context, params dto.MenuUpdateParams) error {
    _, err := s.repo.FindByID(ctx, params.ID.Raw())
    if err != nil {
        return err
    }
    menu := assembler.BuildMenuUpdateCommand(&params)
    _, err = s.repo.Update(ctx, menu)
    return err
}

func (s MenuService) Delete(ctx context.Context, id datatype.SafeUint64) error {
    return s.txManager.Transaction(func(tx *query.Query) error {
        return s.repo.CloneWithQuery(tx).Delete(ctx, id.Raw())
    })
}

func (s MenuService) Move(ctx context.Context, id datatype.SafeUint64, target datatype.SafeUint64) error {
    return s.txManager.Transaction(func(tx *query.Query) error {
        menuRepo := s.repo.CloneWithQuery(tx)
        _, err := menuRepo.FindByIDWithAncestor(ctx, id.Raw(), target.Raw())
        if err == nil {
            return errno.MenuCircular.ToError()
        } else if !errors.Is(err, gorm.ErrRecordNotFound) {
            return err
        }
        
        return menuRepo.Move(ctx, id.Raw(), target.Raw())
    })
}

func (s MenuService) FindByID(ctx context.Context, id datatype.SafeUint64) (*dto.Menu, error) {
    menu, err := s.repo.FindByID(ctx, id.Raw())
    if err != nil {
        return nil, err
    }
    return assembler.BuildMenuDTO(menu), nil
}

func (s MenuService) List(ctx context.Context, params dto.MenuListQueryParams) (*types.PaginatedResult[*dto.Menu], error) {
    data, total, err := s.repo.List(ctx, params)
    if err != nil {
        return nil, err
    }
    return &types.PaginatedResult[*dto.Menu]{
        Pagination: params.Pagination,
        Total:      total,
        List:       assembler.BuildMenuListDTO(data),
    }, nil
}
