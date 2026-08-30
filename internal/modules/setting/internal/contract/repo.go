package contract

import (
	"context"

	"cms/internal/app/setting/model"

	"gorm.io/gen"
)

type Repo interface {
	Add(ctx context.Context, data *model.Setting) error
	Update(ctx context.Context, data *model.Setting) (gen.ResultInfo, error)
	Remove(ctx context.Context, field string) (gen.ResultInfo, error)
	Get(ctx context.Context, field string) (string, error)
	GetAll(ctx context.Context) ([]*model.Setting, error)
}
