package valueobject

import (
    "time"
    
    `cms/internal/modules/contenttype/domain`
    
    "github.com/shopspring/decimal"
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

func (b BoolValue) IsEnumValue(values []any) bool {
    bv := bool(b)
    for _, v := range values {
        if bv == v {
            return true
        }
    }
    return false
}

func (b BoolValue) InRange(_ decimal.Decimal, _ decimal.Decimal) bool {
    return true
}

func (b BoolValue) InTimeRange(_ time.Time, _ time.Time) bool {
    return true
}

func (b BoolValue) WriteTo(s domain.Writable) error {
    bv := bool(b)
    return s.WriteBool(&bv)
}
