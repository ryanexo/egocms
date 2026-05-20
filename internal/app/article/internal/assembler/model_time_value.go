package assembler

import (
    `time`
    
    `cms/internal/app/article/internal/domain`
    `cms/internal/infra/persistence/datatype`
    
    `github.com/shopspring/decimal`
)

type TimeValue time.Time

var _ domain.Value = (*TimeValue)(nil)

func (t TimeValue) IsEmpty() bool {
    return time.Time(t).IsZero()
}

func (t TimeValue) Match(_ string) (bool, error) {
    return true, nil
}

func (t TimeValue) IsValidLen(_ uint64, _ uint64) bool {
    return true
}

func (t TimeValue) IsEnumValue(_ datatype.EnumValues) bool {
    return true
}

func (t TimeValue) InRange(_ decimal.Decimal, _ decimal.Decimal) bool {
    return true
}

func (t TimeValue) InTimeRange(minTime time.Time, maxTime time.Time) bool {
    ts := time.Time(t)
    return ts.After(minTime) && ts.Before(maxTime)
}

func (t TimeValue) Assign(s domain.Scannable) error {
    return s.ScanTime(time.Time(t))
}
