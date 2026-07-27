package valueobject

import (
    "time"
    
    `cms/internal/modules/contenttype/domain`
    
    "github.com/shopspring/decimal"
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

func (t TimeValue) IsEnumValue(values []any) bool {
    tv := time.Time(t)
    for _, v := range values {
        sv, ok := v.(string)
        if !ok {
            return false
        }
        enumTv, err := time.Parse(time.RFC3339, sv)
        if err != nil {
            return false
        }
        if tv.Equal(enumTv) {
            return true
        }
    }
    return false
}

func (t TimeValue) InRange(_ decimal.Decimal, _ decimal.Decimal) bool {
    return true
}

func (t TimeValue) InTimeRange(minTime time.Time, maxTime time.Time) bool {
    ts := time.Time(t)
    return ts.After(minTime) && ts.Before(maxTime)
}

func (t TimeValue) WriteTo(s domain.Writable) error {
    tv := time.Time(t)
    return s.WriteTime(&tv)
}
