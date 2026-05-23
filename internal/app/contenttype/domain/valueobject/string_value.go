package valueobject

import (
    "regexp"
    "time"
    
    `cms/internal/app/contenttype/domain`
    
    "github.com/shopspring/decimal"
)

type StringValue string

var _ domain.Value = (*StringValue)(nil)

func (s StringValue) IsEmpty() bool {
    return s == ""
}

func (s StringValue) Match(pattern string) (bool, error) {
    return regexp.Match(pattern, []byte(s))
}

func (s StringValue) IsValidLen(minLen uint64, maxLen uint64) bool {
    strLen := uint64(len(s))
    return strLen >= minLen && strLen <= maxLen
}

func (s StringValue) IsEnumValue(values map[string]any) bool {
    for _, v := range values {
        if v == string(s) {
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

func (s StringValue) WriteTo(writer domain.Writable) error {
    sv := string(s)
    return writer.WriteString(&sv)
}
