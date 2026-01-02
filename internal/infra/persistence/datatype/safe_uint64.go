package datatype

import (
	"encoding/json"
	"strconv"
	"strings"

	"gorm.io/gorm/schema"
)

type SafeUint64 uint64

var _ json.Marshaler = (*SafeUint64)(nil)
var _ json.Unmarshaler = (*SafeUint64)(nil)
var _ schema.GormDataTypeInterface = (*SafeUint64)(nil)

func (i SafeUint64) MarshalJSON() ([]byte, error) {
	return []byte(`"` + strconv.FormatUint(uint64(i), 10) + `"`), nil
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

func (i SafeUint64) Raw() uint64 {
	return uint64(i)
}

func (i SafeUint64) GormDataType() string {
	return string(schema.Uint)
}
