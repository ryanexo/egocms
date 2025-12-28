package article

import (
    `time`
    
    `dpcms/internal/app/domain/article`
    `dpcms/internal/infra/persistence/customvalue`
    
    `github.com/shopspring/decimal`
)

type TimeValue time.Time

var _ article.Value = (*TimeValue)(nil)

func (t TimeValue) IsEmpty() bool {
    return time.Time(t).IsZero()
}

func (t TimeValue) Match(_ string) (bool, error) {
    return true, nil
}

func (t TimeValue) IsValidLen(_ int, _ int) bool {
    return true
}

func (t TimeValue) IsEnumValue(_ customvalue.EnumValues) bool {
    return true
}

func (t TimeValue) InRange(_ decimal.Decimal, _ decimal.Decimal) bool {
    return true
}

func (t TimeValue) InTimeRange(minTime time.Time, maxTime time.Time) bool {
    ts := time.Time(t)
    return ts.After(minTime) && ts.Before(maxTime)
}

func (t TimeValue) Assign(s article.Scannable) error {
    return s.ScanTime(time.Time(t))
}
