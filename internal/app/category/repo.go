package category

import (
	"context"
	"database/sql"

	"cms/internal/app/category/api"
	"cms/internal/app/category/model"
	"cms/internal/infra/store/closuretable"
	"cms/internal/infra/store/gorm/dbscope"
	"cms/internal/infra/store/gorm/gquery"

	"gorm.io/gen"
)

type Repo struct {
	q       *gquery.Query
	closure closuretable.ClosureTable
}

func NewCategoryRepo(q *gquery.Query, db *sql.DB) (*Repo, error) {
	closure, err := closuretable.NewClosureTable(
		db,
		closuretable.WithTableName(q.CategoryContext.TableName()),
	)
	if err != nil {
		return nil, err
	}
	return &Repo{q: q, closure: closure}, nil
}

func (s *Repo) Create(ctx context.Context, category *model.Category) error {
	err := s.q.Category.WithContext(ctx).Create(category)
	if err != nil {
		return err
	}
	if category.ID > 0 {
		err = s.closure.CreateSubtree(ctx, category.ID, category.ParentID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Repo) Update(ctx context.Context, data *model.Category) (gen.ResultInfo, error) {
	var result gen.ResultInfo
	err := s.q.Transaction(func(tx *gquery.Query) error {
		category := tx.Category
		var err error
		result, err = category.WithContext(ctx).
			Where(category.ID.Eq(data.ID)).
			Select(category.Sequence, category.Name, category.Path).
			Updates(data)
		if err != nil {
			return err
		}

		meta := tx.CategoryMeta
		_, err = meta.WithContext(ctx).
			Where(meta.CategoryID.Eq(data.ID)).
			Select(meta.Title, meta.Keywords, meta.Description, meta.Thumb).
			Updates(&data.Meta)
		return err
	})
	return result, err
}

func (s *Repo) Move(ctx context.Context, fromNode uint64, toNode uint64) error {
	if err := s.closure.Move(ctx, fromNode, toNode); err != nil {
		return err
	}

	category := s.q.Category
	_, err := category.WithContext(ctx).
		Where(category.ID.Eq(fromNode)).
		Update(category.ParentID, toNode)
	return err
}

func (s *Repo) Delete(ctx context.Context, id uint64) error {
	return s.q.Transaction(func(tx *gquery.Query) error {
		_, err := tx.Category.WithContext(ctx).Where(tx.Category.ID.Eq(id)).Delete()
		if err != nil {
			return err
		}
		_, err = tx.CategoryMeta.WithContext(ctx).Where(tx.CategoryMeta.CategoryID.Eq(id)).Delete()
		if err != nil {
			return err
		}
		err = s.closure.DropSubtree(ctx, id)
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *Repo) FindByID(ctx context.Context, id uint64) (*model.Category, error) {
	return s.q.Category.WithContext(ctx).Preload(s.q.Category.Meta).Where(s.q.Category.ID.Eq(id)).First()
}

func (s *Repo) List(ctx context.Context, params *api.CategoryListParams) ([]*model.Category, int64, error) {
	cat := s.q.Category
	meta := s.q.CategoryMeta
	q := cat.WithContext(ctx).Scopes(dbscope.Paginate(params.PageNo, params.PageSize))
	if params.ID != nil {
		q = q.Where(cat.ID.Eq(params.ID.Uint64()))
	}
	if params.ParentID != nil {
		q = q.Where(cat.ParentID.Eq(params.ParentID.Uint64()))
	}
	if params.Name != nil {
		q = q.Where(cat.Name.Like("%" + *params.Name + "%"))
	}
	if params.Keywords != nil || params.Description != nil || params.Title != nil {
		q = q.LeftJoin(meta, cat.ID.EqCol(meta.CategoryID))

		if params.Keywords != nil {
			q = q.Where(meta.Keywords.In(params.Keywords...))
		}
		if params.Title != nil {
			q = q.Where(meta.Title.Like(*params.Title))
		}
		if params.Description != nil {
			q = q.Where(meta.Description.Like(*params.Description))
		}
	}

	total, err := q.Count()
	if err != nil {
		return nil, 0, err
	}

	queryResult, err := q.Preload(cat.Meta).Find()
	if err != nil {
		return nil, 0, err
	}

	return queryResult, total, nil
}
