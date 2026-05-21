package domain

import (
    `cms/internal/app/article/internal/errno`
)

type ModelValue struct {
    value Value
    rules Rules
}

func NewModelValue(rules Rules, value Value) *ModelValue {
    return &ModelValue{rules: rules, value: value}
}

func (v *ModelValue) IsValid() error {
    if !v.rules.IsEnable() {
        return errno.ErrModelDataDisabled.Format(v.rules.FieldKey())
    }
    
    if v.rules.IsRequired() && v.value.IsEmpty() {
        return errno.ErrModelDataMissingValue.Format(v.rules.FieldName())
    }
    
    lenMin, lenMax := v.rules.MinLen(), v.rules.MaxLen()
    if lenMin < lenMax && !v.value.IsValidLen(lenMin, v.rules.MaxLen()) {
        return errno.ErrModelDataInvalidLen.Format(lenMin, v.rules.MaxLen())
    }
    
    vMin, vMax := v.rules.MinValue(), v.rules.MaxValue()
    if vMin.Decimal.LessThan(vMax.Decimal) {
        validMin, validMax := vMin.Valid, vMax.Valid
        if !validMin || !validMax || !v.value.InRange(vMin.Decimal, vMax.Decimal) {
            return errno.ErrModelDataInvalidNumRange.Format(v.rules.MinValue(), v.rules.MaxValue())
        }
    }
    
    if len(v.rules.EnumOptions()) > 0 && !v.value.IsEnumValue(v.rules.EnumOptions()) {
        return errno.ErrModelDataInvalidValue.Format(v.rules.FieldName())
    }
    
    pattern := v.rules.Pattern()
    if len(pattern) > 0 {
        matched, err := v.value.Match(pattern)
        if err != nil {
            return err
        }
        if !matched {
            return errno.ErrModelDataInvalidValue.Format(v.rules.FieldName())
        }
    }
    
    vMinTime, vMaxTime := v.rules.MinTime(), v.rules.MaxTime()
    if !vMinTime.Time.IsZero() || !vMaxTime.Time.IsZero() {
        validMin, validMax := vMinTime.Valid, vMaxTime.Valid
        if !validMin || !validMax || !v.value.InTimeRange(vMinTime.Time, vMaxTime.Time) {
            return errno.ErrModelDataInvalidTimeRange.Format(v.rules.FieldName(), vMinTime, vMaxTime)
        }
    }
    
    return nil
}

func (v *ModelValue) Assign(s Scannable) error {
    return v.value.Assign(s)
}
