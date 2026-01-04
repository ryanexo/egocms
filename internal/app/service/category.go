package service

import (
    "context"
    `errors`
    
    catAssembler `dpcms/internal/app/assembler/category`
    `dpcms/internal/app/dto`
    `dpcms/internal/app/dto/type`
    `dpcms/internal/app/erroz`
    `dpcms/internal/infra`
    `dpcms/internal/infra/persistence/datatype`
    `dpcms/internal/infra/persistence/dbscope`
    `dpcms/internal/infra/persistence/query`
    `dpcms/internal/infra/persistence/repo`
    
    `gorm.io/gorm`
)

type Category struct {
    persist *query.Query
    db      *gorm.DB
}

func NewCategoryService(i *infra.Infra) *Category {
    return &Category{persist: i.Query, db: i.DB}
}

func (srv Category) Create(ctx context.Context, params dto.CategoryCreateParams) (*dto.Category, error) {
    cat, err := catAssembler.ToCategoryCreateCommand(&params)
    if err != nil {
        return nil, err
    }
    err = srv.persist.Transaction(func(tx *query.Query) error {
        catRepo := repo.NewCategoryRepo(tx)
        txErr := catRepo.Create(ctx, cat)
        if txErr != nil {
            return txErr
        }
        txErr = catRepo.CreateSubtree(ctx, cat.ID.Raw(), cat.ParentID.Raw())
        if txErr != nil {
            return txErr
        }
        return nil
    })
    if err != nil {
        return nil, err
    }
    return catAssembler.ToCategoryDTO(cat)
}

func (srv Category) Update(ctx context.Context, params dto.CategoryUpdateParams) error {
    catRepo := repo.NewCategoryRepo(srv.persist)
    _, err := catRepo.FindByID(ctx, params.ID.Raw())
    if err != nil {
        return err
    }
    cat, err := catAssembler.ToCategoryUpdateCommand(&params)
    if err != nil {
        return err
    }
    _, err = catRepo.Update(ctx, cat)
    return err
}

func (srv Category) Delete(ctx context.Context, id datatype.SafeUint64) error {
    return srv.persist.Transaction(func(tx *query.Query) error {
        return repo.NewCategoryRepo(tx).Delete(ctx, id.Raw())
    })
}

func (srv Category) Move(ctx context.Context, id datatype.SafeUint64, target datatype.SafeUint64) error {
    return srv.persist.Transaction(func(tx *query.Query) error {
        catRepo := repo.NewCategoryRepo(tx)
        _, err := catRepo.FindByIDWithAncestor(ctx, id.Raw(), target.Raw())
        if err == nil {
            return erroz.CategoryCircular.ToError()
        } else if !errors.Is(err, gorm.ErrRecordNotFound) {
            return err
        }
        
        return catRepo.Move(ctx, id.Raw(), target.Raw())
    })
}

func (srv Category) FindByID(ctx context.Context, id datatype.SafeUint64) (*dto.Category, error) {
    cat, err := repo.NewCategoryRepo(srv.persist).FindByID(ctx, id.Raw())
    if err != nil {
        return nil, err
    }
    return catAssembler.ToCategoryDTO(cat)
}

func (srv Category) ListRootNodes(ctx context.Context, pageNo int, pageSize int) (*dtotype.PaginatedResult[*dto.Category], error) {
    pid := datatype.SafeUint64(0)
    return srv.List(ctx, dto.CategoryListParams{
        ParentID:   &pid,
        Pagination: dbscope.Pagination{PageNo: pageNo, PageSize: pageSize},
    })
}

func (srv Category) ListNodesByParentID(ctx context.Context, id datatype.SafeUint64, pageSize int, pageNo int) (*dtotype.PaginatedResult[*dto.Category], error) {
    return srv.List(ctx, dto.CategoryListParams{
        ParentID:   &id,
        Pagination: dbscope.Pagination{PageNo: pageNo, PageSize: pageSize},
    })
}

func (srv Category) List(ctx context.Context, params dto.CategoryListParams) (*dtotype.PaginatedResult[*dto.Category], error) {
    data, total, err := repo.NewCategoryRepo(srv.persist).List(ctx, params)
    if err != nil {
        return nil, err
    }
    list, err := catAssembler.ToCategoryListDTO(data)
    if err != nil {
        return nil, err
    }
    
    return &dtotype.PaginatedResult[*dto.Category]{
        Total:    total,
        PageSize: params.PageSize,
        PageNo:   params.PageNo,
        List:     list,
    }, nil
}
