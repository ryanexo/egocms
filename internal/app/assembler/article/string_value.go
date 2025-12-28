package article

import (
    `regexp`
    `time`
    
    `dpcms/internal/app/domain/article`
    `dpcms/internal/infra/persistence/customvalue`
    
    `github.com/shopspring/decimal`
)

type StringValue string

var _ article.Value = (*StringValue)(nil)

func (s StringValue) IsEmpty() bool {
    return s == ""
}

func (s StringValue) Match(pattern string) (bool, error) {
    return regexp.Match(pattern, []byte(s))
}

func (s StringValue) IsValidLen(minLen int, maxLen int) bool {
    strLen := len(s)
    return strLen >= minLen && strLen <= maxLen
}

func (s StringValue) IsEnumValue(values customvalue.EnumValues) bool {
    for _, v := range values {
        if v.Value == string(s) {
            return true
        }
    }
    return false
}

func (s StringValue) InRange(_ decimal.Decimal, _ decimal.Decimal) bool {
    return true
}

func (s StringValue) InTimeRange(_ time.Time, _ time.Time) bool {
    return true
}

func (s StringValue) Assign(scannable article.Scannable) error {
    return scannable.ScanString(string(s))
}
