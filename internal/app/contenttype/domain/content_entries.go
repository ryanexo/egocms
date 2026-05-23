package domain

import (
    `time`
    
    `github.com/shopspring/decimal`
)

type ContentEntries struct {
    fieldKey    string
    typ         int16
    stringValue *string
    numberValue *decimal.Decimal
    timeValue   *time.Time
    boolValue   *bool
}

func NewContentEntries(fieldKey string, typ valueType) ContentEntries {
    return ContentEntries{fieldKey: fieldKey, typ: typ.typ}
}

func (c ContentEntries) FieldKey() string {
    return c.fieldKey
}

func (c ContentEntries) Type() int16 {
    return c.typ
}

func (c ContentEntries) StringValue() *string {
    return c.stringValue
}

func (c ContentEntries) NumberValue() *decimal.Decimal {
    return c.numberValue
}

func (c ContentEntries) TimeValue() *time.Time {
    return c.timeValue
}

func (c ContentEntries) BoolValue() *bool {
    return c.boolValue
}
