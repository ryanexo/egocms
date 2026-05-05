package domain

import (
    `database/sql`
    `time`
    
    `cms/internal/infra/persistence/datatype`
    
    `github.com/shopspring/decimal`
)

type Value interface {
    IsEmpty() bool
    Match(pattern string) (bool, error)
    IsValidLen(minLen uint64, maxLen uint64) bool
    IsEnumValue(values datatype.EnumValues) bool
    InRange(minValue decimal.Decimal, maxValue decimal.Decimal) bool
    InTimeRange(minValue time.Time, maxValue time.Time) bool
    Assign(s Scannable) error
}

type Rules interface {
    FieldKey() string
    FieldName() string
    MinLen() uint64
    MaxLen() uint64
    MinValue() decimal.NullDecimal
    MaxValue() decimal.NullDecimal
    MinTime() sql.NullTime
    MaxTime() sql.NullTime
    Pattern() string
    EnumOptions() datatype.EnumValues
    IsRequired() bool
    IsEnable() bool
}

type Scannable interface {
    ScanString(value any) error
    ScanNumber(value any) error
    ScanBool(value any) error
    ScanTime(value any) error
}

const (
    ValueString = iota
    ValueNumber
    ValueBool
    ValueTime
)
