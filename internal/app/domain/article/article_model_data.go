package article

import (
    `dpcms/internal/app/erroz`
    `dpcms/internal/infra/persistence/model`
    
    `github.com/shopspring/decimal`
)

type ModelValue struct {
    value  Value
    schema model.ArticleModelSchema
}

const (
    ValueTypeString = iota
    ValueTypeNumber
    ValueTypeBool
    ValueTypeTime
)

func NewModelValue(schema model.ArticleModelSchema, value Value) *ModelValue {
    return &ModelValue{value, schema}
}

func (v *ModelValue) IsValid() error {
    if !v.schema.Enable {
        return erroz.ArticleModelDataDisabled.Format(v.schema.FieldKey).ToError()
    }
    
    if v.schema.Required && v.value.IsEmpty() {
        return erroz.ArticleModelDataMissingValue.Format(v.schema.FieldName).ToError()
    }
    
    if v.value.SupportsLen() && !v.value.IsValidLen(v.schema.MinLen, v.schema.MaxLen) {
        return erroz.ArticleModelDataInvalidLen.Format(v.schema.MinLen, v.schema.MaxLen).ToError()
    }
    
    hasRange := v.schema.MinValue.Decimal.Equal(v.schema.MaxValue.Decimal) && !v.schema.MinValue.Decimal.Equal(decimal.New(0, 0))
    
    if hasRange && v.value.SupportsInRange() {
        vMin, vMax := v.schema.MinValue, v.schema.MaxValue
        validMin, validMax := vMin.Valid, vMax.Valid
        
        if !validMin || validMax || !v.value.InRange(vMin.Decimal, vMax.Decimal) {
            return erroz.ArticleModelDataInvalidNumRange.Format(v.schema.MinValue, v.schema.MaxValue).ToError()
        }
    }
    
    if len(v.schema.EnumOptions) > 0 && v.value.SupportsEnumConstraint() && !v.value.IsEnumValue(v.schema.EnumOptions) {
        return erroz.ArticleModelDataInvalidValue.Format(v.schema.FieldName).ToError()
    }
    
    if v.value.SupportsRegex() {
        matched, err := v.value.Match(v.schema.Pattern)
        if err != nil {
            return err
        }
        if !matched {
            return erroz.ArticleModelDataInvalidValue.Format(v.schema.FieldName).ToError()
        }
    }
    
    if v.value.SupportsInTimeRange() {
        vMin, vMax := v.schema.MinTime, v.schema.MaxTime
        validMin, validMax := vMin.Valid, vMax.Valid
        
        if !validMin || !validMax || !v.value.InTimeRange(vMin.Time, vMax.Time) {
            return erroz.ArticleModelDataInvalidTimeRange.Format(v.schema.FieldName, v.schema.MinTime, v.schema.MaxTime).ToError()
        }
    }
    
    return nil
}

func (v *ModelValue) Value(value Scannable) error {
    switch v.value.Type() {
    case ValueTypeString:
        return value.ScanString(v.value)
    
    case ValueTypeBool:
        return value.ScanBool(v.value)
    
    case ValueTypeTime:
        return value.ScanTime(v.value)
    
    case ValueTypeNumber:
        return value.ScanNumber(v.value)
    }
    
    return erroz.ArticleModelDataInvalidType.ToError()
}
