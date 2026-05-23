package domain

import (
    "time"
    
    "github.com/shopspring/decimal"
)

type Value interface {
    IsEmpty() bool
    IsValidLen(minLen uint64, maxLen uint64) bool
    IsEnumValue(values map[string]any) bool
    InRange(minValue decimal.Decimal, maxValue decimal.Decimal) bool
    InTimeRange(minValue time.Time, maxValue time.Time) bool
    Match(pattern string) (bool, error)
    WriteTo(s Writable) error
}

type Writable interface {
    WriteString(value *string) error
    WriteNumber(value *decimal.Decimal) error
    WriteBool(value *bool) error
    WriteTime(value *time.Time) error
}
