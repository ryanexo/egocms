package article

import (
    `dpcms/internal/app/domain/article`
    `dpcms/internal/infra/persistence/model`
)

type ModelData model.ArticleModelData

var _ article.Scannable = (*ModelData)(nil)

func NewModelData(data *model.ArticleModelData) *ModelData {
    return (*ModelData)(data)
}

func (m *ModelData) ScanString(value any) error {
    return m.ValueString.Scan(value)
}

func (m *ModelData) ScanNumber(value any) error {
    return m.ValueNumber.Scan(value)
}

func (m *ModelData) ScanBool(value any) error {
    return m.ValueBool.Scan(value)
}

func (m *ModelData) ScanTime(value any) error {
    return m.ValueTime.Scan(value)
}
