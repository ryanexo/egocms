package article

import (
    `dpcms/internal/infra/persistence/customvalue`
    
    `github.com/shopspring/decimal`
)

type Value interface {
    Type() int16
    SupportsLen() bool
    SupportsInRange() bool
    SupportsEnumConstraint() bool
    SupportsRegex() bool
    IsEmpty() bool
    Match(string) (bool, error)
    IsValidLen(int, int) bool
    IsValidFraction(any) bool
    IsEnumValue(customvalue.EnumValues) bool
    InRange(decimal.NullDecimal, decimal.NullDecimal) bool
}

type Scannable interface {
    ScanString(any) error
    ScanNumber(any) error
    ScanBool(any) error
    ScanFloat(any) error
    ScanTime(any) error
}
