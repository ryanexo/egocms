package valueobject

import (
    "time"
    
    `cms/internal/app/contenttype/domain`
    `cms/internal/app/contenttype/domain/errno`
    
    "github.com/shopspring/decimal"
)

type Value struct {
    v domain.Value
}

func NewValue(typ domain.Schema, value any) (Value, error) {
    if value == nil {
        return Value{NilValue{}}, nil
    }
    
    if domain.ValueBool.Is(typ.Type()) {
        if v, ok := value.(bool); ok {
            return Value{BoolValue(v)}, nil
        }
    } else if domain.ValueNumber.Is(typ.Type()) {
        switch dv := value.(type) {
        case []byte:
            return Value{NumberValue(decimal.RequireFromString(string(dv)))}, nil
        case string:
            n, err := decimal.NewFromString(dv)
            if err != nil {
                return Value{}, err
            }
            return Value{NumberValue(n)}, nil
        case float64:
            return Value{NumberValue(decimal.NewFromFloat(dv))}, nil
        case float32:
            return Value{NumberValue(decimal.NewFromFloat(float64(dv)))}, nil
        case int8:
            return Value{NumberValue(decimal.NewFromInt32(int32(dv)))}, nil
        case int16:
            return Value{NumberValue(decimal.NewFromInt32(int32(dv)))}, nil
        case int32:
            return Value{NumberValue(decimal.NewFromInt32(dv))}, nil
        case int64:
            return Value{NumberValue(decimal.NewFromInt(dv))}, nil
        case uint8:
            return Value{NumberValue(decimal.NewFromUint64(uint64(dv)))}, nil
        case uint16:
            return Value{NumberValue(decimal.NewFromUint64(uint64(dv)))}, nil
        case uint32:
            return Value{NumberValue(decimal.NewFromUint64(uint64(dv)))}, nil
        case uint64:
            return Value{NumberValue(decimal.NewFromUint64(dv))}, nil
        }
    } else if domain.ValueString.Is(typ.Type()) {
        if v, ok := value.(string); ok {
            return Value{StringValue(v)}, nil
        }
    } else if domain.ValueTime.Is(typ.Type()) {
        if v, ok := value.(time.Time); ok {
            return Value{TimeValue(v)}, nil
        }
    }
    return Value{}, errno.ErrInvalidValueType
}

func (v Value) WriteTo(schema domain.Schema, w domain.Writable) error {
    if !schema.Enable() {
        return errno.ErrDisabled.Format(schema.FieldKey())
    }
    
    if schema.Required() && v.v.IsEmpty() {
        return errno.ErrMissingValue.Format(schema.FieldName())
    }
    
    lenMin, lenMax := schema.MinLen(), schema.MaxLen()
    if lenMin < lenMax && !v.v.IsValidLen(lenMin, schema.MaxLen()) {
        return errno.ErrInvalidLen.Format(lenMin, schema.MaxLen())
    }
    
    vMin, vMax := schema.MinValue(), schema.MaxValue()
    if vMin.Decimal.LessThan(vMax.Decimal) {
        validMin, validMax := vMin.Valid, vMax.Valid
        if !validMin || !validMax || !v.v.InRange(vMin.Decimal, vMax.Decimal) {
            return errno.ErrInvalidNumRange.Format(schema.MinValue(), schema.MaxValue())
        }
    }
    
    if len(schema.EnumOptions()) > 0 && !v.v.IsEnumValue(schema.EnumOptions()) {
        return errno.ErrInvalidValue.Format(schema.FieldName())
    }
    
    pattern := schema.Pattern()
    if len(pattern) > 0 {
        matched, err := v.v.Match(pattern)
        if err != nil {
            return err
        }
        if !matched {
            return errno.ErrInvalidValue.Format(schema.FieldName())
        }
    }
    
    vMinTime, vMaxTime := schema.MinTime(), schema.MaxTime()
    if !vMinTime.Time.IsZero() || !vMaxTime.Time.IsZero() {
        validMin, validMax := vMinTime.Valid, vMaxTime.Valid
        if !validMin || !validMax || !v.v.InTimeRange(vMinTime.Time, vMaxTime.Time) {
            return errno.ErrInvalidTimeRange.Format(schema.FieldName(), vMinTime, vMaxTime)
        }
    }
    
    if err := v.v.WriteTo(w); err != nil {
        return err
    }
    
    return nil
}
