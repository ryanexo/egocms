package article

import (
    `time`
    
    `dpcms/internal/infra/persistence/customvalue`
    
    `github.com/shopspring/decimal`
)

type Value interface {
    Type() int16
    SupportsLen() bool
    SupportsInRange() bool
    SupportsEnumConstraint() bool
    SupportsRegex() bool
    SupportsInTimeRange() bool
    IsEmpty() bool
    Match(pattern string) (bool, error)
    IsValidLen(minLen int, maxLen int) bool
    IsEnumValue(values customvalue.EnumValues) bool
    InRange(minValue decimal.Decimal, maxValue decimal.Decimal) bool
    InTimeRange(minValue time.Time, maxValue time.Time) bool
}

type Scannable interface {
    ScanString(value any) error
    ScanNumber(value any) error
    ScanBool(value any) error
    ScanTime(value any) error
}
