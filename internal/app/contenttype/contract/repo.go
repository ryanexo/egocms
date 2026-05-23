package contract

import (
    "context"
    
    "cms/internal/infra/persistence"
    "cms/internal/pkg/datatype"
    
    "cms/internal/app/contenttype/internal/dto"
    "cms/internal/infra/persistence/model"
    
    "gorm.io/gen"
)

type ArticleModelRepo interface {
    persistence.Repository[ArticleModelRepo]
    Create(ctx context.Context, data *model.ContentType) error
    CreateModelJsonData(ctx context.Context, data *model.ContentEntries) error
    CreateModelTypedData(ctx context.Context, data []*model.ContentFieldValues) error
    UpdateModel(ctx context.Context, data *model.ContentType) error
    ReplaceSchema(ctx context.Context, id datatype.SafeUint64, data []*model.ContentTypeSchema) error
    UpdateModelTypedData(ctx context.Context, data []*model.ContentFieldValues) error
    UpdateModelJsonData(ctx context.Context, data *model.ContentEntries) error
    FindByID(ctx context.Context, id datatype.SafeUint64) (*model.ContentType, error)
    FindAllSchema(ctx context.Context, modelID datatype.SafeUint64) ([]*model.ContentTypeSchema, error)
    DeleteModel(ctx context.Context, modelID datatype.SafeUint64) error
    DeleteSchema(ctx context.Context, id datatype.SafeUint64) (gen.ResultInfo, error)
    DeleteAllSchema(ctx context.Context, id datatype.SafeUint64) error
    DeleteArticleData(ctx context.Context, articleID datatype.SafeUint64) error
    List(ctx context.Context, params dto.ArticleModelListParams) ([]*model.ContentType, int64, error)
}
