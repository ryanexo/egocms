package jsontype

import (
    "encoding/json"
    "strconv"
    "strings"
)

type SafeUint64 uint64

var _ json.Marshaler = (*SafeUint64)(nil)
var _ json.Unmarshaler = (*SafeUint64)(nil)

func (i SafeUint64) MarshalJSON() ([]byte, error) {
    return []byte(`"` + i.String() + `"`), nil
}

func (i *SafeUint64) UnmarshalJSON(s []byte) error {
    strVal := strings.Trim(string(s), `"`)
    v, err := strconv.ParseUint(strVal, 10, 64)
    if err != nil {
        return err
    }
    *i = SafeUint64(v)
    return nil
}

func (i SafeUint64) Uint64() uint64 {
    return uint64(i)
}

func (i SafeUint64) String() string {
    return strconv.FormatUint(uint64(i), 10)
}
