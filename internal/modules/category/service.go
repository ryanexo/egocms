package category

import (
    "context"
    "errors"
    
    `cms/internal/infra/store/gorm/gquery`
    `cms/internal/public/apitype`
    `cms/internal/public/jsontype`
    
    contract2 "cms/internal/modules/category/contract"
    "cms/internal/modules/category/internal/assembler"
    "cms/internal/modules/category/internal/dto"
    "cms/internal/modules/category/internal/errno"
    
    "gorm.io/gorm"
)

type CategoryService struct {
    txManager persistence.Transactor
    repo      contract2.CategoryRepo
}

func NewCategoryService(txManager persistence.Transactor, repo contract2.CategoryRepo) *CategoryService {
    return &CategoryService{
        txManager: txManager,
        repo:      repo,
    }
}

func (s CategoryService) Create(ctx context.Context, params dto.CategoryCreateParams) (jsontype.SafeUint64, error) {
    data := assembler.ToCategoryCreateCommand(&params)
    err := s.txManager.Transaction(func(tx *gquery.Query) error {
        catRepo := s.repo.CloneWithQuery(tx)
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
    _, err := s.repo.FindByID(ctx, params.ID.Raw())
    if err != nil {
        return err
    }
    data := assembler.ToCategoryUpdateCommand(&params)
    _, err = s.repo.Update(ctx, data)
    return err
}

func (s CategoryService) Delete(ctx context.Context, id jsontype.SafeUint64) error {
    return s.txManager.Transaction(func(tx *gquery.Query) error {
        return s.repo.CloneWithQuery(tx).Delete(ctx, id.Raw())
    })
}

func (s CategoryService) Move(ctx context.Context, id jsontype.SafeUint64, target jsontype.SafeUint64) error {
    return s.txManager.Transaction(func(tx *gquery.Query) error {
        catRepo := s.repo.CloneWithQuery(tx)
        _, err := catRepo.FindByIDWithAncestor(ctx, id.Raw(), target.Raw())
        if err == nil {
            return errno.ErrCircular.ToError()
        } else if !errors.Is(err, gorm.ErrRecordNotFound) {
            return err
        }
        
        return catRepo.Move(ctx, id.Raw(), target.Raw())
    })
}

func (s CategoryService) FindByID(ctx context.Context, id jsontype.SafeUint64) (*dto.Category, error) {
    data, err := s.repo.FindByID(ctx, id.Raw())
    if err != nil {
        return nil, err
    }
    return assembler.ToCategoryDTO(data), nil
}

func (s CategoryService) ListRootNodes(ctx context.Context, pageNo int, pageSize int) (*apitype.PaginatedResult[*dto.Category], error) {
    pid := jsontype.SafeUint64(0)
    return s.List(ctx, dto.CategoryListParams{
        ParentID:   &pid,
        Pagination: apitype.Pagination{PageNo: pageNo, PageSize: pageSize},
    })
}

func (s CategoryService) ListNodesByParentID(ctx context.Context, id jsontype.SafeUint64, pageSize int, pageNo int) (*apitype.PaginatedResult[*dto.Category], error) {
    return s.List(ctx, dto.CategoryListParams{
        ParentID:   &id,
        Pagination: apitype.Pagination{PageNo: pageNo, PageSize: pageSize},
    })
}

func (s CategoryService) List(ctx context.Context, params dto.CategoryListParams) (*apitype.PaginatedResult[*dto.Category], error) {
    data, total, err := s.repo.List(ctx, params)
    if err != nil {
        return nil, err
    }
    return &apitype.PaginatedResult[*dto.Category]{
        Pagination: params.Pagination,
        Total:      total,
        List:       assembler.ToCategoryListDTO(data),
    }, nil
}
