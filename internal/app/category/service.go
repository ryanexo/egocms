package category

import (
	"context"

	"cms/internal/app/category/api"
	"cms/internal/app/category/model"
	"cms/internal/infra/store/modeltype"
	"cms/internal/public/apitype"
	"cms/internal/public/jsontype"

	"gorm.io/gen"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *model.Category) error
	Update(ctx context.Context, category *model.Category) (gen.ResultInfo, error)
	Move(ctx context.Context, fromNode uint64, toNode uint64) error
	Delete(ctx context.Context, id uint64) error
	FindByID(ctx context.Context, id uint64) (*model.Category, error)
	List(ctx context.Context, params *api.CategoryListParams) ([]*model.Category, int64, error)
}

type Service struct {
	repo CategoryRepository
}

func NewCategoryService(repo CategoryRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, params api.CategoryCreateParams) (jsontype.SafeUint64, error) {
	category, err := buildCategory(params.Name, params.Path, params.Meta)
	if err != nil {
		return 0, err
	}
	category.ParentID = params.ParentID.Uint64()
	category.Sequence = int64(params.Sequence)

	if err = s.repo.Create(ctx, category); err != nil {
		return 0, err
	}
	return jsontype.SafeUint64(category.ID), nil
}

func (s *Service) Update(ctx context.Context, params api.CategoryUpdateParams) error {
	if _, err := s.repo.FindByID(ctx, params.ID.Uint64()); err != nil {
		return err
	}

	category, err := buildCategory(params.Name, params.Path, params.Meta)
	if err != nil {
		return err
	}
	category.Base = modeltype.Base{ID: params.ID.Uint64()}
	category.Sequence = int64(params.Sequence)
	category.Meta.CategoryID = category.ID

	_, err = s.repo.Update(ctx, category)
	return err
}

func (s *Service) Move(ctx context.Context, id, targetID jsontype.SafeUint64) error {
	if _, err := s.repo.FindByID(ctx, id.Uint64()); err != nil {
		return err
	}
	if targetID != 0 {
		if _, err := s.repo.FindByID(ctx, targetID.Uint64()); err != nil {
			return err
		}
	}
	return s.repo.Move(ctx, id.Uint64(), targetID.Uint64())
}

func (s *Service) Delete(ctx context.Context, id jsontype.SafeUint64) error {
	return s.repo.Delete(ctx, id.Uint64())
}

func (s *Service) FindByID(ctx context.Context, id jsontype.SafeUint64) (*model.Category, error) {
	return s.repo.FindByID(ctx, id.Uint64())
}

func (s *Service) ListRootNodes(ctx context.Context, pageNo, pageSize int) (*apitype.PaginatedResult[*model.Category], error) {
	parentID := jsontype.SafeUint64(0)
	return s.List(ctx, api.CategoryListParams{
		Pagination: apitype.Pagination{PageNo: pageNo, PageSize: pageSize},
		ParentID:   &parentID,
	})
}

func (s *Service) ListNodesByParentID(
	ctx context.Context,
	parentID jsontype.SafeUint64,
	pageNo, pageSize int,
) (*apitype.PaginatedResult[*model.Category], error) {
	return s.List(ctx, api.CategoryListParams{
		Pagination: apitype.Pagination{PageNo: pageNo, PageSize: pageSize},
		ParentID:   &parentID,
	})
}

func (s *Service) List(ctx context.Context, params api.CategoryListParams) (*apitype.PaginatedResult[*model.Category], error) {
	list, total, err := s.repo.List(ctx, &params)
	if err != nil {
		return nil, err
	}
	return &apitype.PaginatedResult[*model.Category]{
		Pagination: params.Pagination,
		Total:      total,
		List:       list,
	}, nil
}

func buildCategory(name, categoryPath string, metaParams api.CategoryMetaParams) (*model.Category, error) {
	category := &model.Category{}
	if err := category.SetName(name); err != nil {
		return nil, err
	}
	if err := category.SetPath(categoryPath); err != nil {
		return nil, err
	}
	if err := category.Meta.SetTitle(metaParams.Title); err != nil {
		return nil, err
	}
	if err := category.Meta.SetKeywords(metaParams.Keywords); err != nil {
		return nil, err
	}
	if err := category.Meta.SetDescription(metaParams.Description); err != nil {
		return nil, err
	}
	if err := category.Meta.SetThumb(metaParams.Thumb); err != nil {
		return nil, err
	}
	return category, nil
}
