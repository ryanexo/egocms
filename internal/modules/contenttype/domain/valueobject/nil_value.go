package valueobject

import (
    "time"
    
    `cms/internal/modules/contenttype/domain`
    
    "github.com/shopspring/decimal"
)

type NilValue struct{}

var _ domain.Value = (*NilValue)(nil)

func (n NilValue) IsEmpty() bool {
    return true
}

func (n NilValue) Match(_ string) (bool, error) {
    return true, nil
}

func (n NilValue) IsValidLen(_ uint64, _ uint64) bool {
    return true
}

func (n NilValue) IsEnumValue(_ []any) bool {
    return true
}

func (n NilValue) InRange(_ decimal.Decimal, _ decimal.Decimal) bool {
    return true
}

func (n NilValue) InTimeRange(_ time.Time, _ time.Time) bool {
    return true
}

func (n NilValue) WriteTo(writer domain.Writable) error {
    if err := writer.WriteString(nil); err != nil {
        return err
    }
    if err := writer.WriteNumber(nil); err != nil {
        return err
    }
    if err := writer.WriteBool(nil); err != nil {
        return err
    }
    if err := writer.WriteTime(nil); err != nil {
        return err
    }
    return nil
}
