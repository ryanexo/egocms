package dto

import (
    `dpcms/internal/infra/persistence/datatype`
)

type ResourceID struct {
    ID datatype.SafeUint64 `validate:"required" json:"id" swaggertype:"string"`
}
