package assembler

import (
    `time`
    
    `cms/internal/domain/article/internal/domain`
    `cms/internal/infra/persistence/datatype`
    
    `github.com/shopspring/decimal`
)

type BoolValue bool

var _ domain.Value = (*BoolValue)(nil)

func (b BoolValue) IsEmpty() bool {
    return false
}

func (b BoolValue) Match(_ string) (bool, error) {
    return true, nil
}

func (b BoolValue) IsValidLen(_ uint64, _ uint64) bool {
    return true
}

func (b BoolValue) IsEnumValue(_ datatype.EnumValues) bool {
    return true
}

func (b BoolValue) InRange(_ decimal.Decimal, _ decimal.Decimal) bool {
    return true
}

func (b BoolValue) InTimeRange(_ time.Time, _ time.Time) bool {
    return true
}

func (b BoolValue) Assign(s domain.Scannable) error {
    return s.ScanBool(bool(b))
}
