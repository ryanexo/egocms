package article

import (
    `context`
    
    `cms/internal/domain/articlemodel/internal/dto`
    `cms/internal/infra/persistence/contract`
    `cms/internal/infra/persistence/datatype`
    `cms/internal/infra/persistence/model`
    
    `gorm.io/gen`
)

type ArticleModelRepo interface {
    contract.Repository[ArticleModelRepo]
    Create(ctx context.Context, data *model.ArticleModel) error
    CreateModelJsonData(ctx context.Context, data *model.ArticleModelJsonData) error
    CreateModelTypedData(ctx context.Context, data []*model.ArticleModelData) error
    UpdateModel(ctx context.Context, data *model.ArticleModel) error
    ReplaceSchema(ctx context.Context, id datatype.SafeUint64, data []*model.ArticleModelSchema) error
    UpdateModelTypedData(ctx context.Context, data []*model.ArticleModelData) error
    UpdateModelJsonData(ctx context.Context, data *model.ArticleModelJsonData) error
    FindByID(ctx context.Context, id datatype.SafeUint64) (*model.ArticleModel, error)
    FindAllSchema(ctx context.Context, modelID datatype.SafeUint64) ([]*model.ArticleModelSchema, error)
    DeleteModel(ctx context.Context, modelID datatype.SafeUint64) error
    DeleteSchema(ctx context.Context, id datatype.SafeUint64) (gen.ResultInfo, error)
    DeleteAllSchema(ctx context.Context, id datatype.SafeUint64) error
    DeleteArticleData(ctx context.Context, articleID datatype.SafeUint64) error
    List(ctx context.Context, params dto.ArticleModelListParams) ([]*model.ArticleModel, int64, error)
}
