package erroz

import (
    `strings`
    
    `dpcms/internal/erroz/internal/errcode`
    `dpcms/internal/erroz/internal/errmod`
    `dpcms/internal/erroz/internal/errtype`
)

func Code(module errmod.Code, etype errtype.Code, code errcode.Code) string {
    return strings.Join([]string{
        module.String(),
        etype.String(),
        code.String(6),
    }, "-")
}
