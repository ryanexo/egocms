package customvalue

import (
    `database/sql/driver`
    `fmt`
    
    `github.com/bytedance/sonic`
)

type EnumValue struct {
    Label string `json:"label"`
    Value string `json:"value"`
}

type EnumValues []EnumValue

func (v *EnumValues) Scan(value any) error {
    if value == nil || value == "" {
        *v = []EnumValue{}
        return nil
    }
    
    strVal, ok := value.(string)
    if !ok {
        return fmt.Errorf("invalid type for EnumValues: %T", value)
    }
    
    enumValues := make(EnumValues, 0)
    err := sonic.UnmarshalString(strVal, &enumValues)
    if err != nil {
        return err
    }
    *v = enumValues
    
    return nil
}

func (v EnumValues) Value() (driver.Value, error) {
    if v == nil {
        return nil, nil
    }
    return sonic.MarshalString(v)
}
