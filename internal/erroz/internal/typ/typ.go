package typ

import `dpcms/internal/erroz/internal/format`

type Typ struct {
    Code int64
}

func (s Typ) String() string {
    return format.Code(s.Code, 3)
}
