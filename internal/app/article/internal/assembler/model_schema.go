package assembler

import (
    `database/sql`
    
    `cms/internal/app/article/internal/domain`
    `cms/internal/infra/persist/datatype`
    `cms/internal/infra/persist/model`
    
    `github.com/shopspring/decimal`
)

type ModelSchema struct {
    data *model.ArticleModelSchema
}

var _ domain.Rules = (*ModelSchema)(nil)

func (m ModelSchema) FieldKey() string {
    return m.data.FieldKey
}

func (m ModelSchema) FieldName() string {
    return m.data.FieldName
}

func (m ModelSchema) MinLen() uint64 {
    return m.data.MinLen.Raw()
}

func (m ModelSchema) MaxLen() uint64 {
    return m.data.MaxLen.Raw()
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

func (m ModelSchema) EnumOptions() datatype.EnumValues {
    return m.data.EnumOptions
}

func (m ModelSchema) IsRequired() bool {
    return m.data.Required == 1
}

func (m ModelSchema) IsEnable() bool {
    return m.data.Enable == 1
}

func NewRules(data *model.ArticleModelSchema) domain.Rules {
    return ModelSchema{data}
}
