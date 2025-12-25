package article

import (
    `dpcms/internal/app/erroz`
    `dpcms/internal/infra/persistence/model`
)

type ArticleModelData struct {
    value  Value
    schema model.ArticleModelSchema
}

const (
    ValueTypeString = iota
    ValueTypeNumber
    ValueTypeBool
    ValueTypeFloat
    ValueTypeTime
)

func NewArticleModelData(schema model.ArticleModelSchema, value Value) *ArticleModelData {
    return &ArticleModelData{value: value, schema: schema}
}

func (v *ArticleModelData) IsValid() error {
    if !v.schema.Enable {
        return erroz.ArticleModelDataDisabled.Format(v.schema.FieldKey).ToError()
    }
    if v.schema.Required && v.value.IsEmpty() {
        return erroz.ArticleModelDataMissingValue.Format(v.schema.FieldName).ToError()
    }
    if v.value.SupportsLen() && !v.value.IsValidLen(v.schema.MinLen, v.schema.MaxLen) {
        return erroz.ArticleModelDataInvalidLen.Format(v.schema.MinLen, v.schema.MaxLen).ToError()
    }
    if v.value.SupportsInRange() && !v.value.InRange(v.schema.MinValue, v.schema.MaxValue) {
        return erroz.ArticleModelDataInvalidRange.Format(v.schema.MinValue, v.schema.MaxValue).ToError()
    }
    if v.value.SupportsEnumConstraint() && !v.value.IsEnumValue(v.schema.EnumOptions) {
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
    
    return nil
}

func (v *ArticleModelData) Value(value Scannable) error {
    switch v.value.Type() {
    case ValueTypeString:
        return value.ScanString(v.value)
    
    case ValueTypeBool:
        return value.ScanBool(v.value)
    
    case ValueTypeFloat:
        return value.ScanFloat(v.value)
    
    case ValueTypeTime:
        return value.ScanTime(v.value)
    
    case ValueTypeNumber:
        return value.ScanNumber(v.value)
    }
    
    return erroz.ArticleModelDataInvalidType.ToError()
}
