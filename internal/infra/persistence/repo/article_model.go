package repo

import (
    `context`
    
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/infra/persistence/query`
)

type ArticleModel struct {
    persist *query.Query
}

func NewArticleModel(persist *query.Query) *ArticleModel {
    return &ArticleModel{persist}
}

func (r *ArticleModel) FindSchema(ctx context.Context, modelId uint64) ([]*model.ArticleModelSchema, error) {
    schemaModel := r.persist.ArticleModelSchema
    return schemaModel.WithContext(ctx).Where(schemaModel.ModelId.Eq(modelId)).Find()
}

func (r *ArticleModel) Create(ctx context.Context, data *model.ArticleModel) error {
    return r.persist.ArticleModel.WithContext(ctx).Create(data)
}

func (r *ArticleModel) CreateModelJsonData(ctx context.Context, data *model.ArticleModelJsonData) error {
    return r.persist.ArticleModelJsonData.WithContext(ctx).Create(data)
}

func (r *ArticleModel) CreateModelTypedData(ctx context.Context, data []*model.ArticleModelData) error {
    return r.persist.ArticleModelData.WithContext(ctx).Create(data...)
}

func (r *ArticleModel) UpdateModelTypedData(ctx context.Context, data []*model.ArticleModelData) error {
    if len(data) == 0 {
        return nil
    }
    typedModel := r.persist.ArticleModelData
    for _, item := range data {
        _, err := typedModel.WithContext(ctx).Where(
            typedModel.ModelId.Eq(item.ModelId),
            typedModel.ArticleId.Eq(item.ArticleId),
            typedModel.FieldKey.Eq(item.FieldKey),
        ).Updates(item)
        
        if err != nil {
            return err
        }
    }
    return nil
}

func (r *ArticleModel) UpdateModelJsonData(ctx context.Context, data *model.ArticleModelJsonData) error {
    jsonModel := r.persist.ArticleModelJsonData
    _, err := jsonModel.WithContext(ctx).Where(
        jsonModel.ModelId.Eq(data.ModelId),
        jsonModel.ArticleId.Eq(data.ArticleId),
    ).Update(jsonModel.Data, data.Data)
    return err
}
