package jsontype

import (
    "encoding/json"
    "strconv"
    "strings"
)

type SafeInt64 int64

var _ json.Marshaler = (*SafeInt64)(nil)
var _ json.Unmarshaler = (*SafeInt64)(nil)

func (i SafeInt64) MarshalJSON() ([]byte, error) {
    return []byte(`"` + i.String() + `"`), nil
}

func (i *SafeInt64) UnmarshalJSON(s []byte) error {
    strVal := strings.Trim(string(s), `"`)
    v, err := strconv.ParseInt(strVal, 10, 64)
    if err != nil {
        return err
    }
    *i = SafeInt64(v)
    return nil
}

func (i SafeInt64) Int64() int64 {
    return int64(i)
}

func (i SafeInt64) String() string {
    return strconv.FormatInt(int64(i), 10)
}
