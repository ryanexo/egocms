package setting

import (
	"context"

	"cms/internal/app/setting/model"
	"cms/internal/infra/store/gorm/gquery"

	"gorm.io/gen"
)

type Repo struct {
	q *gquery.Query
}

func NewSettingRepo(q *gquery.Query) *Repo {
	return &Repo{q: q}
}

func (s *Repo) Create(ctx context.Context, data *model.Setting) error {
	return s.q.Setting.WithContext(ctx).Create(data)
}

func (s *Repo) Add(ctx context.Context, data *model.Setting) error {
	return s.Create(ctx, data)
}

func (s *Repo) Update(ctx context.Context, data *model.Setting) (gen.ResultInfo, error) {
	setting := s.q.Setting
	return setting.WithContext(ctx).
		Where(setting.Field.Eq(data.Field)).
		Update(setting.Value, data.Value)
}

func (s *Repo) Delete(ctx context.Context, field string) (gen.ResultInfo, error) {
	setting := s.q.Setting
	return setting.WithContext(ctx).Where(setting.Field.Eq(field)).Delete()
}

func (s *Repo) Remove(ctx context.Context, field string) (gen.ResultInfo, error) {
	return s.Delete(ctx, field)
}

func (s *Repo) FindByField(ctx context.Context, field string) (*model.Setting, error) {
	setting := s.q.Setting
	return setting.WithContext(ctx).Where(setting.Field.Eq(field)).First()
}

func (s *Repo) Get(ctx context.Context, field string) (string, error) {
	data, err := s.FindByField(ctx, field)
	if err != nil {
		return "", err
	}
	return data.Value, nil
}

func (s *Repo) List(ctx context.Context) ([]*model.Setting, error) {
	return s.q.Setting.WithContext(ctx).Find()
}

func (s *Repo) GetAll(ctx context.Context) ([]*model.Setting, error) {
	return s.List(ctx)
}
