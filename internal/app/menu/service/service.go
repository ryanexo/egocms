package service

import (
    `context`
    `errors`
    
    `cms/internal/app/menu/internal/assembler`
    `cms/internal/app/menu/internal/dto`
    `cms/internal/app/menu/internal/errno`
    `cms/internal/infra/persist/contract`
    `cms/internal/infra/persist/datatype`
    `cms/internal/infra/persist/model`
    `cms/internal/infra/persist/query`
    `cms/internal/util/types`
    
    `gorm.io/gorm`
)

type MenuService struct {
    txManager contract.Transactor
    repo      contract.MenuRepo
}

func NewMenuService(txManager contract.Transactor, repo contract.MenuRepo) *MenuService {
    return &MenuService{
        txManager: txManager,
        repo:      repo,
    }
}

func (s MenuService) Create(ctx context.Context, params dto.MenuCreateParams) (*model.Menu, error) {
    menu := assembler.ToMenuCreateCommand(&params)
    err := s.txManager.Transaction(func(tx *query.Query) error {
        menuRepo := s.repo.CloneWithQuery(tx)
        txErr := menuRepo.Create(ctx, menu)
        if txErr != nil {
            return txErr
        }
        txErr = menuRepo.CreateSubtree(ctx, menu.ID, menu.ParentID)
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
    _, err := s.repo.FindByID(ctx, params.ID)
    if err != nil {
        return err
    }
    menu := assembler.ToMenuUpdateCommand(&params)
    _, err = s.repo.Update(ctx, menu)
    return err
}

func (s MenuService) Delete(ctx context.Context, id datatype.SafeUint64) error {
    return s.txManager.Transaction(func(tx *query.Query) error {
        return s.repo.CloneWithQuery(tx).Delete(ctx, id)
    })
}

func (s MenuService) Move(ctx context.Context, id datatype.SafeUint64, target datatype.SafeUint64) error {
    return s.txManager.Transaction(func(tx *query.Query) error {
        menuRepo := s.repo.CloneWithQuery(tx)
        _, err := menuRepo.FindByIDWithAncestor(ctx, id, target)
        if err == nil {
            return errno.MenuCircular.ToError()
        } else if !errors.Is(err, gorm.ErrRecordNotFound) {
            return err
        }
        
        return menuRepo.Move(ctx, id, target)
    })
}

func (s MenuService) FindByID(ctx context.Context, id datatype.SafeUint64) (*dto.Menu, error) {
    menu, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    return assembler.ToMenuDTO(menu), nil
}

func (s MenuService) List(ctx context.Context, params dto.MenuListQueryParams) (*types.PaginatedResult[*dto.Menu], error) {
    data, total, err := s.repo.List(ctx, params)
    if err != nil {
        return nil, err
    }
    return &types.PaginatedResult[*dto.Menu]{
        Pagination: params.Pagination,
        Total:      total,
        List:       assembler.ToMenuListDTO(data),
    }, nil
}
