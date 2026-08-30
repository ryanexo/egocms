package assembler

import (
    `cms/internal/modules/contenttype/domain`
    `cms/internal/public/model`
)

type ModelData model.ContentFieldValues

var _ domain.Writable = (*ModelData)(nil)

func (m *ModelData) SetString(value any) error {
    s := (*model.ContentFieldValues)(m)
    return s.StringValue.Scan(value)
}

func (m *ModelData) SetNumber(value any) error {
    s := (*model.ContentFieldValues)(m)
    return s.NumberValue.Scan(value)
}

func (m *ModelData) SetBool(value any) error {
    s := (*model.ContentFieldValues)(m)
    return s.BoolValue.Scan(value)
}

func (m *ModelData) SetTime(value any) error {
    s := (*model.ContentFieldValues)(m)
    return s.TimeValue.Scan(value)
}

func (m *ModelData) Model() *model.ContentFieldValues {
    return (*model.ContentFieldValues)(m)
}

func NewModelData(data *model.ContentFieldValues) *ModelData {
    return (*ModelData)(data)
}
