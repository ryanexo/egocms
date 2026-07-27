package valueobject

import (
    "time"
    
    `cms/internal/modules/contenttype/domain`
    `cms/internal/modules/contenttype/internal/errno`
    
    "github.com/shopspring/decimal"
)

type Value struct {
    v domain.Value
}

func NewValue(typ domain.Rule, value any) (Value, error) {
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

func (v Value) WriteTo(rule domain.Rule, w domain.Writable) error {
    if !rule.Enable() {
        return errno.ErrDisabled.Format(rule.FieldKey())
    }
    
    if rule.Required() && v.v.IsEmpty() {
        return errno.ErrMissingValue.Format(rule.FieldName())
    }
    
    lenMinPtr, lenMaxPtr := rule.MinLen(), rule.MaxLen()
    if lenMinPtr != nil && lenMaxPtr != nil {
        lenMin, lenMax := *lenMinPtr, *lenMaxPtr
        if lenMin < lenMax && !v.v.IsValidLen(lenMin, lenMax) {
            return errno.ErrInvalidLen.Format(lenMin, lenMax)
        }
    }
    
    vMinPtr, vMaxPtr := rule.MinValue(), rule.MaxValue()
    if vMinPtr != nil && vMaxPtr != nil {
        vMin, vMax := *vMinPtr, *vMaxPtr
        if !v.v.InRange(vMin, vMax) {
            return errno.ErrInvalidNumRange.Format(rule.FieldName(), rule.MinValue(), rule.MaxValue())
        }
    }
    
    if len(rule.EnumOptions()) > 0 && !v.v.IsEnumValue(rule.EnumOptions()) {
        return errno.ErrInvalidValue.Format(rule.FieldName())
    }
    
    pattern := rule.Pattern()
    if pattern != nil && len(*pattern) > 0 {
        matched, err := v.v.Match(*pattern)
        if err != nil {
            return err
        }
        if !matched {
            return errno.ErrInvalidValue.Format(rule.FieldName())
        }
    }
    
    vMinTimePtr, vMaxTimePtr := rule.MinTime(), rule.MaxTime()
    if vMinTimePtr != nil && vMaxTimePtr != nil {
        vMinTime, vMaxTime := *vMinTimePtr, *vMaxTimePtr
        if !vMinTime.IsZero() || !vMaxTime.IsZero() || !v.v.InTimeRange(vMinTime, vMaxTime) {
            return errno.ErrInvalidTimeRange.Format(rule.FieldName(), vMinTime, vMaxTime)
        }
    }
    
    if err := v.v.WriteTo(w); err != nil {
        return err
    }
    
    return nil
}
