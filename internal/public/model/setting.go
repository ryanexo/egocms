package model

import (
    `cms/internal/infra/store/modeltype`
)

type Setting struct {
    modeltype.Base
    Field string
    Value string
}
