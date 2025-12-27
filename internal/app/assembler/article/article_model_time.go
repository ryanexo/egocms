package article

import (
    `time`
    
    `dpcms/internal/app/domain/article`
    `dpcms/internal/infra/persistence/customvalue`
    
    `github.com/shopspring/decimal`
)

type TimeValue time.Time

var _ article.Value = (*TimeValue)(nil)

func (t TimeValue) Type() int16 {
    return article.ValueTypeTime
}

func (t TimeValue) SupportsLen() bool {
    return false
}

func (t TimeValue) SupportsInRange() bool {
    return true
}

func (t TimeValue) SupportsInTimeRange() bool {
    return true
}

func (t TimeValue) SupportsEnumConstraint() bool {
    return false
}

func (t TimeValue) SupportsRegex() bool {
    return false
}

func (t TimeValue) IsEmpty() bool {
    return time.Time(t).IsZero()
}

func (t TimeValue) Match(_ string) (bool, error) {
    return false, nil
}

func (t TimeValue) IsValidLen(_ int, _ int) bool {
    return false
}

func (t TimeValue) IsEnumValue(_ customvalue.EnumValues) bool {
    return false
}

func (t TimeValue) InRange(_ decimal.Decimal, _ decimal.Decimal) bool {
    return false
}

func (t TimeValue) InTimeRange(minTime time.Time, maxTime time.Time) bool {
    ts := time.Time(t)
    return ts.After(minTime) && ts.Before(maxTime)
}
