package article

import (
    `dpcms/internal/app/domain/article`
    `dpcms/internal/infra/persistence/model`
)

type ModelData model.ArticleModelData

var _ article.Scannable = (*ModelData)(nil)

func (m *ModelData) ScanString(value any) error {
    s := (*model.ArticleModelData)(m)
    return s.ValueString.Scan(value)
}

func (m *ModelData) ScanNumber(value any) error {
    s := (*model.ArticleModelData)(m)
    return s.ValueNumber.Scan(value)
}

func (m *ModelData) ScanBool(value any) error {
    s := (*model.ArticleModelData)(m)
    return s.ValueBool.Scan(value)
}

func (m *ModelData) ScanTime(value any) error {
    s := (*model.ArticleModelData)(m)
    return s.ValueTime.Scan(value)
}

func NewModelData(defaultValue model.ArticleModelData) (*ModelData, *model.ArticleModelData) {
    persist := &model.ArticleModelData{
        ModelId:   defaultValue.ModelId,
        ArticleId: defaultValue.ArticleId,
        FieldKey:  defaultValue.FieldKey,
        Type:      defaultValue.Type,
    }
    return (*ModelData)(persist), persist
}
