package errtype

import `cms/internal/erroz/internal/format`

type ErrType int

func (s ErrType) String() string {
    return format.Code(int64(s), 3)
}

const (
    OK ErrType = iota
    Unknown
    Auth
    NotFound
    InvalidState
    Parameter
    Conflict
)
