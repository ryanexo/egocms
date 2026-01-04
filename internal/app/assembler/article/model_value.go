package article

import (
    `errors`
    `time`
    
    `dpcms/internal/app/domain/article`
    `dpcms/internal/infra/persistence/datatype`
    
    `github.com/shopspring/decimal`
)

func NewValue(vType int16, value any) (article.Value, error) {
    if value == nil {
        return NilValue{}, nil
    }
    
    switch vType {
    case article.ValueBool:
        if v, ok := value.(bool); ok {
            return BoolValue(v), nil
        }
    
    case article.ValueNumber:
        switch dv := value.(type) {
        case []byte:
            return NumberValue(decimal.RequireFromString(string(dv))), nil
        case string:
            n, err := decimal.NewFromString(dv)
            if err != nil {
                return nil, err
            }
            return NumberValue(n), nil
        case float64:
            return NumberValue(decimal.NewFromFloat(dv)), nil
        case float32:
            return NumberValue(decimal.NewFromFloat(float64(dv))), nil
        case int8:
            return NumberValue(decimal.NewFromInt32(int32(dv))), nil
        case int16:
            return NumberValue(decimal.NewFromInt32(int32(dv))), nil
        case int32:
            return NumberValue(decimal.NewFromInt32(dv)), nil
        case int64:
            return NumberValue(decimal.NewFromInt(dv)), nil
        case uint8:
            return NumberValue(decimal.NewFromUint64(uint64(dv))), nil
        case uint16:
            return NumberValue(decimal.NewFromUint64(uint64(dv))), nil
        case uint32:
            return NumberValue(decimal.NewFromUint64(uint64(dv))), nil
        case uint64:
            return NumberValue(decimal.NewFromUint64(dv)), nil
        case datatype.SafeUint64:
            return NumberValue(decimal.NewFromUint64(dv.Raw())), nil
        }
    
    case article.ValueString:
        if v, ok := value.(string); ok {
            return StringValue(v), nil
        }
    
    case article.ValueTime:
        if v, ok := value.(time.Time); ok {
            return TimeValue(v), nil
        }
    }
    
    return nil, errors.New("unsupported type")
}
