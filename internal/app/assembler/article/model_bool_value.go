package article

import (
    `time`
    
    `dpcms/internal/app/domain/article`
    `dpcms/internal/infra/persistence/customvalue`
    
    `github.com/shopspring/decimal`
)

type BoolValue bool

var _ article.Value = (*BoolValue)(nil)

func (b BoolValue) IsEmpty() bool {
    return false
}

func (b BoolValue) Match(_ string) (bool, error) {
    return true, nil
}

func (b BoolValue) IsValidLen(_ int, _ int) bool {
    return true
}

func (b BoolValue) IsEnumValue(_ customvalue.EnumValues) bool {
    return true
}

func (b BoolValue) InRange(_ decimal.Decimal, _ decimal.Decimal) bool {
    return true
}

func (b BoolValue) InTimeRange(_ time.Time, _ time.Time) bool {
    return true
}

func (b BoolValue) Assign(s article.Scannable) error {
    return s.ScanBool(bool(b))
}
