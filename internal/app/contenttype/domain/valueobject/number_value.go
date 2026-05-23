package valueobject

import (
    "regexp"
    "time"
    
    `cms/internal/app/contenttype/domain`
    
    "github.com/shopspring/decimal"
)

type NumberValue decimal.Decimal

var _ domain.Value = (*NumberValue)(nil)

func (n NumberValue) IsEmpty() bool {
    return false
}

func (n NumberValue) Match(pattern string) (bool, error) {
    d := decimal.Decimal(n)
    strValue := d.String()
    return regexp.Match(pattern, []byte(strValue))
}

func (n NumberValue) IsValidLen(_ uint64, _ uint64) bool {
    return true
}

func (n NumberValue) IsEnumValue(values map[string]any) bool {
    d := decimal.Decimal(n)
    for _, v := range values {
        if d.Float64() == v {
            return true
        }
    }
    return false
}

func (n NumberValue) InRange(minValue decimal.Decimal, maxValue decimal.Decimal) bool {
    d := decimal.Decimal(n)
    return d.GreaterThanOrEqual(minValue) && d.LessThanOrEqual(maxValue)
}

func (n NumberValue) InTimeRange(_ time.Time, _ time.Time) bool {
    return true
}

func (n NumberValue) WriteTo(writer domain.Writable) error {
    dv := decimal.Decimal(n)
    return writer.WriteNumber(&dv)
}
