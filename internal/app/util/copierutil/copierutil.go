package copierutil

import (
    `strconv`
    
    `dpcms/internal/infra/persistence/customvalue`
    
    `github.com/jinzhu/copier`
)

func WithUint64ToString() copier.TypeConverter {
    return copier.TypeConverter{
        SrcType: uint64(0),
        DstType: "",
        Fn: func(src any) (dst any, err error) {
            return strconv.FormatUint(src.(uint64), 10), nil
        },
    }
}

func WithStringToUint64() copier.TypeConverter {
    return copier.TypeConverter{
        SrcType: customvalue.Uint64String(""),
        DstType: uint64(0),
        Fn: func(src any) (dst any, err error) {
            return strconv.ParseUint(src.(string), 10, 64)
        },
    }
}
