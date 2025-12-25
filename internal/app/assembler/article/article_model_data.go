package article

import (
    `dpcms/internal/app/domain/article`
    `dpcms/internal/app/erroz`
    `dpcms/internal/infra/persistence/model`
)

type ArticleModelData struct{}

func NewArticleModelData(schema model.ArticleModelSchema, value any) (*model.ArticleModelData, error) {
    data := model.ArticleModelData{ModelId: schema.ModelId, FieldKey: schema.FieldKey, Type: schema.Type}
    
    convRules := []struct {
        Src int16
        Dst interface{ Scan(any) error }
    }{
        {
            Src: article.ValueTypeString,
            Dst: &data.ValueString,
        },
        {
            Src: article.ValueTypeNumber,
            Dst: &data.ValueNumber,
        },
        {
            Src: article.ValueTypeBool,
            Dst: &data.ValueBool,
        },
        {
            Src: article.ValueTypeFloat,
            Dst: &data.ValueFloat,
        },
        {
            Src: article.ValueTypeTime,
            Dst: &data.ValueTime,
        },
    }
    
    for _, rule := range convRules {
        if rule.Src != data.Type {
            continue
        }
        if err := rule.Dst.Scan(value); err == nil {
            return &ArticleModelData{schema, data}, nil
        }
    }
    
    return nil, erroz.ArticleModelDataInvalidType.Format(schema.FieldName).ToError()
}
