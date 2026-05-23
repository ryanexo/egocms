package domain

import (
    `regexp`
    `time`
    
    `github.com/shopspring/decimal`
)

type Schema struct {
    fieldKey    string
    fieldName   string
    description string
    sequence    int64
    typ         int16
    minLen      *uint64
    maxLen      *uint64
    minValue    *decimal.Decimal
    maxValue    *decimal.Decimal
    minTime     *time.Time
    maxTime     *time.Time
    pattern     *string
    enumOptions map[string]any
    required    bool
    enable      bool
    visible     bool
}

func NewSchema(fieldKey string, fieldName string) Schema {
    return Schema{
        fieldKey:  fieldKey,
        fieldName: fieldName,
        required:  false,
        enable:    true,
    }
}

func (m Schema) FieldKey() string {
    return m.fieldKey
}

func (m Schema) FieldName() string {
    return m.fieldName
}

func (m Schema) Description() string {
    return m.description
}

func (m Schema) Sequence() int64 {
    return m.sequence
}

func (m Schema) Type() int16 {
    return m.typ
}

func (m Schema) MinLen() *uint64 {
    return m.minLen
}

func (m Schema) MaxLen() *uint64 {
    return m.maxLen
}

func (m Schema) MinValue() *decimal.Decimal {
    return m.minValue
}

func (m Schema) MaxValue() *decimal.Decimal {
    return m.maxValue
}

func (m Schema) MinTime() *time.Time {
    return m.minTime
}

func (m Schema) MaxTime() *time.Time {
    return m.maxTime
}

func (m Schema) Pattern() *string {
    return m.pattern
}

func (m Schema) EnumOptions() map[string]any {
    return m.enumOptions
}

func (m Schema) Required() bool {
    return m.required
}

func (m Schema) Enable() bool {
    return m.enable
}

func (m Schema) SetTextLength(minLen uint64, maxLen uint64) {
    if maxLen < minLen {
        minLen, maxLen = maxLen, minLen
    }
    
    m.minLen = &minLen
    m.maxLen = &maxLen
}

func (m Schema) SetNumberRange(minValue uint64, maxValue uint64) {
    if maxValue < minValue {
        minValue, maxValue = maxValue, minValue
    }
    
    vMin, vMax := decimal.NewFromUint64(minValue), decimal.NewFromUint64(maxValue)
    m.minValue = &vMin
    m.maxValue = &vMax
}

func (m Schema) SetTimeRange(minTime time.Time, maxTime time.Time) {
    if maxTime.Before(minTime) {
        minTime, maxTime = maxTime, minTime
    }
    
    m.minTime = &minTime
    m.maxTime = &maxTime
}

func (m Schema) SetDescription(description string) {
    m.description = description
}

func (m Schema) SetRequired(required bool) {
    m.required = required
}

func (m Schema) SetEnable(enable bool) {
    m.enable = enable
}

func (m Schema) SetPattern(pattern string) error {
    _, err := regexp.Compile(pattern)
    if err != nil {
        return err
    }
    m.pattern = &pattern
    return nil
}

func (m Schema) SetEnumOptions(options map[string]any) {
    m.enumOptions = options
}

func (m Schema) SetSequence(sequence int64) {
    m.sequence = sequence
}

func (m Schema) SetType(vt valueType) {
    m.typ = vt.typ
}

func (m Schema) SetVisible(visible bool) {
    m.enable = visible
}
