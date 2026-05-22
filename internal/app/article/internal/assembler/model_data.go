package assembler

import (
    `cms/internal/app/article/internal/domain`
    `cms/internal/infra/persistence/model`
)

type ModelData model.ArticleModelData

var _ domain.Setter = (*ModelData)(nil)

func (m *ModelData) SetString(value any) error {
    s := (*model.ArticleModelData)(m)
    return s.ValueString.Scan(value)
}

func (m *ModelData) SetNumber(value any) error {
    s := (*model.ArticleModelData)(m)
    return s.ValueNumber.Scan(value)
}

func (m *ModelData) SetBool(value any) error {
    s := (*model.ArticleModelData)(m)
    return s.ValueBool.Scan(value)
}

func (m *ModelData) SetTime(value any) error {
    s := (*model.ArticleModelData)(m)
    return s.ValueTime.Scan(value)
}

func (m *ModelData) Model() *model.ArticleModelData {
    return (*model.ArticleModelData)(m)
}

func NewModelData(data *model.ArticleModelData) *ModelData {
    return (*ModelData)(data)
}
