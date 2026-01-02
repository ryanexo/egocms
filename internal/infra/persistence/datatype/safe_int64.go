package datatype

import (
	"encoding/json"
	"strconv"
	"strings"

	"gorm.io/gorm/schema"
)

type SafeInt64 int64

var _ json.Marshaler = (*SafeInt64)(nil)
var _ json.Unmarshaler = (*SafeInt64)(nil)
var _ schema.GormDataTypeInterface = (*SafeInt64)(nil)

func (i SafeInt64) MarshalJSON() ([]byte, error) {
	return []byte(`"` + strconv.FormatInt(int64(i), 10) + `"`), nil
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

func (i SafeInt64) Raw() int64 {
	return int64(i)
}

func (i SafeInt64) GormDataType() string {
	return string(schema.Int)
}
