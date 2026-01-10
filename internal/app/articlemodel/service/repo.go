package service

import (
    `context`
    
    `dpcms/internal/app/articlemodel/internal/dto`
    `dpcms/internal/infra/persistence/datatype`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/infra/persistence/contract`
    
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
    FindAllSchema(ctx context.Context, modelId datatype.SafeUint64) ([]*model.ArticleModelSchema, error)
    DeleteModel(ctx context.Context, modelId datatype.SafeUint64) error
    DeleteSchema(ctx context.Context, id datatype.SafeUint64) (gen.ResultInfo, error)
    DeleteAllSchema(ctx context.Context, id datatype.SafeUint64) error
    DeleteArticleData(ctx context.Context, articleId datatype.SafeUint64) error
    List(ctx context.Context, params dto.ArticleModelListParams) ([]*model.ArticleModel, int64, error)
}
