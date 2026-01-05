package datatype

import (
    `database/sql/driver`
    `errors`
    
    `github.com/bytedance/sonic`
)

type JSONMap map[string]interface{}

func (j *JSONMap) Scan(value interface{}) error {
    if value == nil {
        *j = nil
        return nil
    }
    
    if v, ok := value.([]byte); ok {
        return sonic.Unmarshal(v, j)
    }
    
    if v, ok := value.(string); ok {
        return sonic.UnmarshalString(v, j)
    }
    
    return errors.New("unsupported Scan type for JSONMap")
}

func (j JSONMap) Value() (driver.Value, error) {
    return sonic.Marshal(j)
}
