package article

import (
    `time`
    
    `dpcms/internal/app/domain/article`
    `dpcms/internal/infra/persistence/datatype`
    
    `github.com/shopspring/decimal`
)

type NilValue struct{}

var _ article.Value = (*NilValue)(nil)

func (n NilValue) IsEmpty() bool {
    return true
}

func (n NilValue) Match(_ string) (bool, error) {
    return true, nil
}

func (n NilValue) IsValidLen(_ uint64, _ uint64) bool {
    return true
}

func (n NilValue) IsEnumValue(_ datatype.EnumValues) bool {
    return true
}

func (n NilValue) InRange(_ decimal.Decimal, _ decimal.Decimal) bool {
    return true
}

func (n NilValue) InTimeRange(_ time.Time, _ time.Time) bool {
    return true
}

func (n NilValue) Assign(_ article.Scannable) error {
    return nil
}
