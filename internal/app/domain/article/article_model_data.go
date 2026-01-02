package article

import (
    `dpcms/internal/app/erroz`
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
        return erroz.ArticleModelDataDisabled.Format(v.rules.FieldKey()).ToError()
    }
    
    if v.rules.IsRequired() && v.value.IsEmpty() {
        return erroz.ArticleModelDataMissingValue.Format(v.rules.FieldName()).ToError()
    }
    
    lenMin, lenMax := v.rules.MinLen(), v.rules.MaxLen()
    if lenMin < lenMax && !v.value.IsValidLen(lenMin, v.rules.MaxLen()) {
        return erroz.ArticleModelDataInvalidLen.Format(lenMin, v.rules.MaxLen()).ToError()
    }
    
    vMin, vMax := v.rules.MinValue(), v.rules.MaxValue()
    if vMin.Decimal.LessThan(vMax.Decimal) {
        validMin, validMax := vMin.Valid, vMax.Valid
        if !validMin || !validMax || !v.value.InRange(vMin.Decimal, vMax.Decimal) {
            return erroz.ArticleModelDataInvalidNumRange.Format(v.rules.MinValue(), v.rules.MaxValue()).ToError()
        }
    }
    
    if len(v.rules.EnumOptions()) > 0 && !v.value.IsEnumValue(v.rules.EnumOptions()) {
        return erroz.ArticleModelDataInvalidValue.Format(v.rules.FieldName()).ToError()
    }
    
    pattern := v.rules.Pattern()
    if len(pattern) > 0 {
        matched, err := v.value.Match(pattern)
        if err != nil {
            return err
        }
        if !matched {
            return erroz.ArticleModelDataInvalidValue.Format(v.rules.FieldName()).ToError()
        }
    }
    
    vMinTime, vMaxTime := v.rules.MinTime(), v.rules.MaxTime()
    if !vMinTime.Time.IsZero() || !vMaxTime.Time.IsZero() {
        validMin, validMax := vMinTime.Valid, vMaxTime.Valid
        if !validMin || !validMax || !v.value.InTimeRange(vMinTime.Time, vMaxTime.Time) {
            return erroz.ArticleModelDataInvalidTimeRange.Format(v.rules.FieldName(), vMinTime, vMaxTime).ToError()
        }
    }
    
    return nil
}

func (v *ModelValue) Assign(s Scannable) error {
    return v.value.Assign(s)
}
