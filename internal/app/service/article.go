package service

import (
    `context`
    `database/sql`
    
    `dpcms/internal/app/domain/article`
    `dpcms/internal/app/dto`
    `dpcms/internal/app/erroz`
    `dpcms/internal/app/repo`
    `dpcms/internal/app/util/copierutil`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/infra/persistence/query`
    
    `github.com/jinzhu/copier`
)

type Article struct {
    query *query.Query
}

func NewArticle(query *query.Query) *Article {
    return &Article{query}
}

func (s Article) Create(ctx context.Context, params dto.ArticleCreateParams) (*dto.ArticleDetail, error) {
    artModel := &model.Article{}
    err := copier.CopyWithOption(&artModel, &params, copier.Option{
        Converters: []copier.TypeConverter{copierutil.WithStringToUint64()},
    })
    if err != nil {
        return nil, err
    }
    
    err = s.query.Transaction(func(tx *query.Query) error {
        return repo.NewArticle(ctx, tx).Create(artModel)
    })
    if err != nil {
        return nil, err
    }
}

func (s Article) transformContentModelDataWithVerify(ctx context.Context, modelId uint64, data map[string]any) ([]*model.ArticleModelData, error) {
    defDao := s.query.ArticleModelSchema
    defs, err := defDao.WithContext(ctx).Where(defDao.ModelId.Eq(modelId)).Find()
    if err != nil {
        return nil, err
    }
    
    result := make([]*model.ArticleModelData, 0, len(defs))
    for _, def := range defs {
        value, found := data[def.FieldKey]
        if !found && def.Required.Bool {
            return nil, erroz.ArticleModelDataMissingValue.Format(def.FieldName).ToError()
        }
        
        if rawVal, ok := value.(string); ok && def.Type == article.ValueTypeString {
            result = append(result, &model.ArticleModelData{
                ModelId:     modelId,
                FieldKey:    def.FieldKey,
                ValueString: sql.NullString{String: rawVal, Valid: true},
            })
            continue
        }
    }
    
    return nil, nil
}
