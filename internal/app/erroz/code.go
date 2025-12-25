package erroz

import (
    `strings`
    
    `dpcms/internal/app/erroz/internal/errcode`
    `dpcms/internal/app/erroz/internal/errmod`
    `dpcms/internal/app/erroz/internal/errtype`
)

func Code(module errmod.Code, etype errtype.Code, code errcode.Code) string {
    return strings.Join([]string{
        module.String(),
        etype.String(),
        code.String(6),
    }, "-")
}
