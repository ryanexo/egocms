package article

import (
    `database/sql`
    `time`
    
    `dpcms/internal/infra/persistence/customvalue`
    
    `github.com/shopspring/decimal`
)

type Value interface {
    IsEmpty() bool
    Match(pattern string) (bool, error)
    IsValidLen(minLen int, maxLen int) bool
    IsEnumValue(values customvalue.EnumValues) bool
    InRange(minValue decimal.Decimal, maxValue decimal.Decimal) bool
    InTimeRange(minValue time.Time, maxValue time.Time) bool
    Assign(s Scannable) error
}

type Rules interface {
    FieldKey() string
    FieldName() string
    MinLen() int
    MaxLen() int
    MinValue() decimal.NullDecimal
    MaxValue() decimal.NullDecimal
    MinTime() sql.NullTime
    MaxTime() sql.NullTime
    Pattern() string
    EnumOptions() customvalue.EnumValues
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
