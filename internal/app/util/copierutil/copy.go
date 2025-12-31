package copierutil

import `github.com/jinzhu/copier`

func CopyWithIdConverter(toValue any, fromValue any) error {
    return copier.CopyWithOption(
        toValue,
        fromValue,
        copier.Option{
            Converters: []copier.TypeConverter{WithUint64ToString(), WithStringToUint64()},
        },
    )
}
