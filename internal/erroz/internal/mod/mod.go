package mod

import `dpcms/internal/erroz/internal/format`

type Module struct {
    Code int64
}

func (s Module) String() string {
    return format.Code(s.Code, 3)
}
