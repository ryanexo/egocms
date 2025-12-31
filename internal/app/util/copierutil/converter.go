package copierutil

import (
    `errors`
    
    `dpcms/internal/infra/persistence/customvalue`
    
    `github.com/jinzhu/copier`
)

func WithUint64ToString() copier.TypeConverter {
    return copier.TypeConverter{
        SrcType: uint64(0),
        DstType: customvalue.Uint64String(""),
        Fn: func(src any) (dst any, err error) {
            typedValue, ok := src.(uint64)
            if !ok {
                return "", errors.New("invalid type for uint64")
            }
            return customvalue.NewUint64String(typedValue), nil
        },
    }
}

func WithStringToUint64() copier.TypeConverter {
    return copier.TypeConverter{
        SrcType: customvalue.Uint64String(""),
        DstType: uint64(0),
        Fn: func(src any) (dst any, err error) {
            typedValue, ok := src.(customvalue.Uint64String)
            if !ok {
                return 0, errors.New("invalid type for Uint64String")
            }
            return typedValue.Uint64()
        },
    }
}
