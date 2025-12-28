package article

import (
    `database/sql`
    
    `dpcms/internal/app/domain/article`
    `dpcms/internal/infra/persistence/customvalue`
    `dpcms/internal/infra/persistence/model`
    
    `github.com/shopspring/decimal`
)

type ModelSchema struct {
    data *model.ArticleModelSchema
}

var _ article.Rules = (*ModelSchema)(nil)

func (m ModelSchema) FieldKey() string {
    return m.data.FieldKey
}

func (m ModelSchema) FieldName() string {
    return m.data.FieldName
}

func (m ModelSchema) MinLen() int {
    return m.data.MinLen
}

func (m ModelSchema) MaxLen() int {
    return m.data.MaxLen
}

func (m ModelSchema) MinValue() decimal.NullDecimal {
    return m.data.MinValue
}

func (m ModelSchema) MaxValue() decimal.NullDecimal {
    return m.data.MaxValue
}

func (m ModelSchema) MinTime() sql.NullTime {
    return m.data.MinTime
}

func (m ModelSchema) MaxTime() sql.NullTime {
    return m.data.MaxTime
}

func (m ModelSchema) Pattern() string {
    return m.data.Pattern
}

func (m ModelSchema) EnumOptions() customvalue.EnumValues {
    return m.data.EnumOptions
}

func (m ModelSchema) IsRequired() bool {
    return m.data.Required
}

func (m ModelSchema) IsEnable() bool {
    return m.data.Enable
}

func NewRules(data *model.ArticleModelSchema) article.Rules {
    return ModelSchema{data}
}
