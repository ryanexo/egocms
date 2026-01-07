package errtype

import (
    `dpcms/internal/erroz/internal/errcode`
)

type Code int64

func (e Code) String() string {
    return errcode.Code(e).String(6)
}

const (
    OK Code = iota
    Unknown
    Auth
    NotFound
    InvalidState
    Parameter
    Conflict
)
