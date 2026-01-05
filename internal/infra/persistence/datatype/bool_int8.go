package datatype

import `encoding/json`

type BoolInt8 int8

var _ json.Marshaler = (*BoolInt8)(nil)
var _ json.Unmarshaler = (*BoolInt8)(nil)

func (i BoolInt8) MarshalJSON() ([]byte, error) {
    if i == 0 {
        return []byte{'f', 'a', 'l', 's', 'e'}, nil
    }
    return []byte{'t', 'r', 'u', 'e'}, nil
}

func (i *BoolInt8) UnmarshalJSON(s []byte) error {
    if len(s) == 4 && s[0] == 't' && s[1] == 'r' && s[2] == 'u' && s[3] == 'e' {
        *i = 1
    } else {
        *i = 0
    }
    return nil
}

func (i BoolInt8) Raw() int8 {
    return int8(i)
}
