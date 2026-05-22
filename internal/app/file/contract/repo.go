package contract

import (
	"cms/internal/infra/persistence"
	"cms/internal/pkg/datatype"
	"context"

	"cms/internal/infra/persistence/model"

	"gorm.io/gen"
)

type FileRepo interface {
	persistence.Repository[FileRepo]
	Create(ctx context.Context, data *model.File) error
	Delete(ctx context.Context, id datatype.SafeUint64) (gen.ResultInfo, error)
	FindByID(ctx context.Context, id datatype.SafeUint64) (*model.File, error)
	FindByPath(ctx context.Context, path string) (*model.File, error)
}
