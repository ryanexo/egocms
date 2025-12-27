package article

import (
    `time`
    
    `dpcms/internal/app/domain/article`
    `dpcms/internal/infra/persistence/customvalue`
    
    `github.com/shopspring/decimal`
)

type BoolValue bool

var _ article.Value = (*BoolValue)(nil)

func (b BoolValue) Type() int16 {
    return article.ValueTypeBool
}

func (b BoolValue) SupportsLen() bool {
    return false
}

func (b BoolValue) SupportsInRange() bool {
    return false
}

func (b BoolValue) SupportsInTimeRange() bool {
    return false
}

func (b BoolValue) SupportsEnumConstraint() bool {
    return false
}

func (b BoolValue) SupportsRegex() bool {
    return false
}

func (b BoolValue) IsEmpty() bool {
    return false
}

func (b BoolValue) Match(_ string) (bool, error) {
    return false, nil
}

func (b BoolValue) IsValidLen(_ int, _ int) bool {
    return false
}

func (b BoolValue) IsEnumValue(_ customvalue.EnumValues) bool {
    return false
}

func (b BoolValue) InRange(_ decimal.Decimal, _ decimal.Decimal) bool {
    return false
}

func (b BoolValue) InTimeRange(_ time.Time, _ time.Time) bool {
    return false
}
