package service

import (
    "context"
    `errors`
    
    `dpcms/internal/app/category/errno`
    `dpcms/internal/app/category/internal/assembler`
    `dpcms/internal/app/category/internal/dto`
    `dpcms/internal/app/category/repo`
    `dpcms/internal/infra/persistence/datatype`
    `dpcms/internal/infra/persistence/query`
    `dpcms/internal/types`
    
    `gorm.io/gorm`
)

type CategoryService struct {
    Query *query.Query
}

func (s CategoryService) Create(ctx context.Context, params dto.CategoryCreateParams) (datatype.SafeUint64, error) {
    data := assembler.ToCategoryCreateCommand(&params)
    err := s.Query.Transaction(func(tx *query.Query) error {
        catRepo := repo.NewRepository(tx)
        txErr := catRepo.Create(ctx, data)
        if txErr != nil {
            return txErr
        }
        txErr = catRepo.CreateSubtree(ctx, data.ID.Raw(), data.ParentID.Raw())
        if txErr != nil {
            return txErr
        }
        return nil
    })
    if err != nil {
        return 0, err
    }
    return data.ID, nil
}

func (s CategoryService) Update(ctx context.Context, params dto.CategoryUpdateParams) error {
    catRepo := repo.NewRepository(s.Query)
    _, err := catRepo.FindByID(ctx, params.ID.Raw())
    if err != nil {
        return err
    }
    data := assembler.ToCategoryUpdateCommand(&params)
    _, err = catRepo.Update(ctx, data)
    return err
}

func (s CategoryService) Delete(ctx context.Context, id datatype.SafeUint64) error {
    return s.Query.Transaction(func(tx *query.Query) error {
        return repo.NewRepository(tx).Delete(ctx, id.Raw())
    })
}

func (s CategoryService) Move(ctx context.Context, id datatype.SafeUint64, target datatype.SafeUint64) error {
    return s.Query.Transaction(func(tx *query.Query) error {
        catRepo := repo.NewRepository(tx)
        _, err := catRepo.FindByIDWithAncestor(ctx, id.Raw(), target.Raw())
        if err == nil {
            return errno.CategoryCircular.ToError()
        } else if !errors.Is(err, gorm.ErrRecordNotFound) {
            return err
        }
        
        return catRepo.Move(ctx, id.Raw(), target.Raw())
    })
}

func (s CategoryService) FindByID(ctx context.Context, id datatype.SafeUint64) (*dto.Category, error) {
    data, err := repo.NewRepository(s.Query).FindByID(ctx, id.Raw())
    if err != nil {
        return nil, err
    }
    return assembler.ToCategoryDTO(data), nil
}

func (s CategoryService) ListRootNodes(ctx context.Context, pageNo int, pageSize int) (*types.PaginatedResult[*dto.Category], error) {
    pid := datatype.SafeUint64(0)
    return s.List(ctx, dto.CategoryListParams{
        ParentID:   &pid,
        Pagination: types.Pagination{PageNo: pageNo, PageSize: pageSize},
    })
}

func (s CategoryService) ListNodesByParentID(ctx context.Context, id datatype.SafeUint64, pageSize int, pageNo int) (*types.PaginatedResult[*dto.Category], error) {
    return s.List(ctx, dto.CategoryListParams{
        ParentID:   &id,
        Pagination: types.Pagination{PageNo: pageNo, PageSize: pageSize},
    })
}

func (s CategoryService) List(ctx context.Context, params dto.CategoryListParams) (*types.PaginatedResult[*dto.Category], error) {
    data, total, err := repo.NewRepository(s.Query).List(ctx, params)
    if err != nil {
        return nil, err
    }
    return &types.PaginatedResult[*dto.Category]{
        Pagination: params.Pagination,
        Total:      total,
        List:       assembler.ToCategoryListDTO(data),
    }, nil
}
