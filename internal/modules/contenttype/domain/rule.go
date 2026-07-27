package domain

import (
    `regexp`
    `time`
    
    `github.com/shopspring/decimal`
)

type Rule struct {
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
    enumOptions []map[string]any
    required    bool
    enable      bool
    visible     bool
}

func NewRule(fieldKey string, fieldName string) Rule {
    return Rule{
        fieldKey:  fieldKey,
        fieldName: fieldName,
        required:  false,
        enable:    true,
    }
}

func (m Rule) FieldKey() string {
    return m.fieldKey
}

func (m Rule) FieldName() string {
    return m.fieldName
}

func (m Rule) Description() string {
    return m.description
}

func (m Rule) Sequence() int64 {
    return m.sequence
}

func (m Rule) Type() int16 {
    return m.typ
}

func (m Rule) MinLen() *uint64 {
    return m.minLen
}

func (m Rule) MaxLen() *uint64 {
    return m.maxLen
}

func (m Rule) MinValue() *decimal.Decimal {
    return m.minValue
}

func (m Rule) MaxValue() *decimal.Decimal {
    return m.maxValue
}

func (m Rule) MinTime() *time.Time {
    return m.minTime
}

func (m Rule) MaxTime() *time.Time {
    return m.maxTime
}

func (m Rule) Pattern() *string {
    return m.pattern
}

func (m Rule) EnumOptions() []any {
    exists := make(map[any]struct{})
    result := make([]any, 0)
    for _, v := range m.enumOptions {
        _, found := exists[v]
        if !found {
            result = append(result, v)
        }
    }
    return result
}

func (m Rule) Required() bool {
    return m.required
}

func (m Rule) Enable() bool {
    return m.enable
}

func (m Rule) SetTextLength(minLen uint64, maxLen uint64) {
    if maxLen < minLen {
        minLen, maxLen = maxLen, minLen
    }
    
    m.minLen = &minLen
    m.maxLen = &maxLen
}

func (m Rule) SetNumberRange(minValue uint64, maxValue uint64) {
    if maxValue < minValue {
        minValue, maxValue = maxValue, minValue
    }
    
    vMin, vMax := decimal.NewFromUint64(minValue), decimal.NewFromUint64(maxValue)
    m.minValue = &vMin
    m.maxValue = &vMax
}

func (m Rule) SetTimeRange(minTime time.Time, maxTime time.Time) {
    if maxTime.Before(minTime) {
        minTime, maxTime = maxTime, minTime
    }
    
    m.minTime = &minTime
    m.maxTime = &maxTime
}

func (m Rule) SetDescription(description string) {
    m.description = description
}

func (m Rule) SetRequired(required bool) {
    m.required = required
}

func (m Rule) SetEnable(enable bool) {
    m.enable = enable
}

func (m Rule) SetPattern(pattern string) error {
    _, err := regexp.Compile(pattern)
    if err != nil {
        return err
    }
    m.pattern = &pattern
    return nil
}

func (m Rule) SetEnumOptions(options []map[string]any) {
    m.enumOptions = options
}

func (m Rule) SetSequence(sequence int64) {
    m.sequence = sequence
}

func (m Rule) SetType(vt valueType) {
    m.typ = vt.v
}

func (m Rule) SetVisible(visible bool) {
    m.enable = visible
}
