package article

import (
    `regexp`
    `time`
    
    `dpcms/internal/app/domain/article`
    `dpcms/internal/infra/persistence/customvalue`
    
    `github.com/shopspring/decimal`
)

type NumberValue decimal.Decimal

var _ article.Value = (*NumberValue)(nil)

func (n NumberValue) IsEmpty() bool {
    return false
}

func (n NumberValue) Match(pattern string) (bool, error) {
    d := decimal.Decimal(n)
    strValue := d.String()
    return regexp.Match(pattern, []byte(strValue))
}

func (n NumberValue) IsValidLen(_ int, _ int) bool {
    return true
}

func (n NumberValue) IsEnumValue(values customvalue.EnumValues) bool {
    d := decimal.Decimal(n)
    for _, v := range values {
        if d.String() == v.Value {
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

func (n NumberValue) Assign(s article.Scannable) error {
    return s.ScanNumber(decimal.Decimal(n))
}
