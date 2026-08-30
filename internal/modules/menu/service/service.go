package service

import (
    "context"
    "errors"
    
    `cms/internal/infra/store/gorm/gquery`
    `cms/internal/public/apitype`
    `cms/internal/public/jsontype`
    `cms/internal/public/model`
    
    contract2 "cms/internal/modules/menu/contract"
    "cms/internal/modules/menu/internal/assembler"
    "cms/internal/modules/menu/internal/dto"
    "cms/internal/modules/menu/internal/errno"
    
    "gorm.io/gorm"
)

type MenuService struct {
    txManager persistence.Transactor
    repo      contract2.MenuRepo
}

func NewMenuService(txManager persistence.Transactor, repo contract2.MenuRepo) *MenuService {
    return &MenuService{
        txManager: txManager,
        repo:      repo,
    }
}

func (s MenuService) Create(ctx context.Context, params dto.MenuCreateParams) (*model.Menu, error) {
    menu := assembler.ToMenuCreateCommand(&params)
    err := s.txManager.Transaction(func(tx *gquery.Query) error {
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

func (s MenuService) Delete(ctx context.Context, id jsontype.SafeUint64) error {
    return s.txManager.Transaction(func(tx *gquery.Query) error {
        return s.repo.CloneWithQuery(tx).Delete(ctx, id)
    })
}

func (s MenuService) Move(ctx context.Context, id jsontype.SafeUint64, target jsontype.SafeUint64) error {
    return s.txManager.Transaction(func(tx *gquery.Query) error {
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

func (s MenuService) FindByID(ctx context.Context, id jsontype.SafeUint64) (*dto.Menu, error) {
    menu, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    return assembler.ToMenuDTO(menu), nil
}

func (s MenuService) List(ctx context.Context, params dto.MenuListQueryParams) (*apitype.PaginatedResult[*dto.Menu], error) {
    data, total, err := s.repo.List(ctx, params)
    if err != nil {
        return nil, err
    }
    return &apitype.PaginatedResult[*dto.Menu]{
        Pagination: params.Pagination,
        Total:      total,
        List:       assembler.ToMenuListDTO(data),
    }, nil
}
