package erroz

import (
    `strings`
    
    `dpcms/internal/erroz/internal/format`
    `dpcms/internal/erroz/type`
)

func Code(module string, typ errtype.ErrType, code int64) string {
    return strings.Join([]string{
        strings.ToUpper(module),
        typ.String(),
        format.Code(code, 6),
    }, "-")
}
