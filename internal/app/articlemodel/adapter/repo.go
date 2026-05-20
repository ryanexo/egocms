package adapter

import (
    "context"
    
    `cms/internal/app/articlemodel/contract`
    `cms/internal/app/articlemodel/internal/dto`
    "cms/internal/infra/persistence/datatype"
    "cms/internal/infra/persistence/dbscope"
    "cms/internal/infra/persistence/model"
    "cms/internal/infra/persistence/query"
    
    "gorm.io/gen"
)

type articleModelRepo struct {
    query *query.Query
}

func NewArticleModelRepo(persist *query.Query) contract.ArticleModelRepo {
    return &articleModelRepo{persist}
}

func (r *articleModelRepo) CloneWithQuery(q *query.Query) contract.ArticleModelRepo {
    return NewArticleModelRepo(q)
}

func (r *articleModelRepo) Create(ctx context.Context, data *model.ArticleModel) error {
    return r.query.ArticleModel.WithContext(ctx).Create(data)
}
func (r *articleModelRepo) CreateModelJsonData(ctx context.Context, data *model.ArticleModelJsonData) error {
    return r.query.ArticleModelJsonData.WithContext(ctx).Create(data)
}

func (r *articleModelRepo) CreateModelTypedData(ctx context.Context, data []*model.ArticleModelData) error {
    return r.query.ArticleModelData.WithContext(ctx).Create(data...)
}

func (r *articleModelRepo) UpdateModel(ctx context.Context, data *model.ArticleModel) error {
    m := r.query.ArticleModel
    _, err := m.WithContext(ctx).Where(m.ID.Eq(data.ID.Raw())).Updates(data)
    return err
}

func (r *articleModelRepo) ReplaceSchema(ctx context.Context, id datatype.SafeUint64, data []*model.ArticleModelSchema) error {
    _, err := r.DeleteSchema(ctx, id)
    if err != nil {
        return err
    }
    m := r.query.ArticleModelSchema
    return m.WithContext(ctx).CreateInBatches(data, 500)
}

func (r *articleModelRepo) UpdateModelTypedData(ctx context.Context, data []*model.ArticleModelData) error {
    if len(data) == 0 {
        return nil
    }
    typedModel := r.query.ArticleModelData
    for _, item := range data {
        _, err := typedModel.WithContext(ctx).Where(
            typedModel.ModelID.Eq(item.ModelID.Raw()),
            typedModel.ArticleID.Eq(item.ArticleID.Raw()),
            typedModel.FieldKey.Eq(item.FieldKey),
        ).Updates(item)
        
        if err != nil {
            return err
        }
    }
    return nil
}

func (r *articleModelRepo) UpdateModelJsonData(ctx context.Context, data *model.ArticleModelJsonData) error {
    jsonModel := r.query.ArticleModelJsonData
    _, err := jsonModel.WithContext(ctx).Where(
        jsonModel.ModelID.Eq(data.ModelID.Raw()),
        jsonModel.ArticleID.Eq(data.ArticleID.Raw()),
    ).Update(jsonModel.Data, data.Data)
    return err
}

func (r *articleModelRepo) FindByID(ctx context.Context, id datatype.SafeUint64) (*model.ArticleModel, error) {
    return r.query.ArticleModel.WithContext(ctx).Where(r.query.ArticleModel.ID.Eq(id.Raw())).First()
}

func (r *articleModelRepo) FindAllSchema(ctx context.Context, modelID datatype.SafeUint64) ([]*model.ArticleModelSchema, error) {
    schemaModel := r.query.ArticleModelSchema
    return schemaModel.WithContext(ctx).Where(schemaModel.ModelID.Eq(modelID.Raw())).Find()
}

func (r *articleModelRepo) DeleteModel(ctx context.Context, modelID datatype.SafeUint64) error {
    _, err := r.query.ArticleModel.WithContext(ctx).Unscoped().Where(r.query.ArticleModel.ID.Eq(modelID.Raw())).Delete()
    return err
}

func (r *articleModelRepo) DeleteSchema(ctx context.Context, id datatype.SafeUint64) (gen.ResultInfo, error) {
    m := r.query.ArticleModelSchema
    return m.WithContext(ctx).Unscoped().Where(m.ModelID.Eq(id.Raw())).Delete()
}

func (r *articleModelRepo) DeleteAllSchema(ctx context.Context, id datatype.SafeUint64) error {
    m := r.query.ArticleModelSchema
    _, err := m.WithContext(ctx).Unscoped().Where(m.ModelID.Eq(id.Raw())).Delete()
    return err
}

func (r *articleModelRepo) DeleteArticleData(ctx context.Context, articleID datatype.SafeUint64) error {
    m := r.query.ArticleModelData
    _, err := m.WithContext(ctx).Unscoped().Where(m.ArticleID.Eq(articleID.Raw())).Delete()
    return err
}

func (r *articleModelRepo) List(ctx context.Context, params dto.ArticleModelListParams) ([]*model.ArticleModel, int64, error) {
    m := r.query.ArticleModel
    q := m.WithContext(ctx).Scopes(dbscope.Paginate(params.PageNo, params.PageSize))
    if params.Name != nil {
        q = q.Where(m.Name.Like("%" + *params.Name + "%"))
    }
    count, err := q.Count()
    if err != nil {
        return nil, 0, err
    }
    result, err := q.Find()
    if err != nil {
        return nil, 0, err
    }
    return result, count, nil
}
