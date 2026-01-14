package erroz

import (
    `strings`
    
    `cms/internal/erroz/internal/format`
    `cms/internal/erroz/type`
)

func Code(module string, typ errtype.ErrType, code int64) string {
    return strings.Join([]string{
        strings.ToUpper(module),
        typ.String(),
        format.Code(code, 6),
    }, "-")
}
