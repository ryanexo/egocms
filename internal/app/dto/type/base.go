package dtotype

import (
    `time`
    
    `dpcms/internal/infra/persistence/datatype`
)

type Base struct {
    ID        datatype.SafeUint64 `json:"id"`
    CreatedAt time.Time           `json:"createdAt"`
    UpdatedAt time.Time           `json:"updatedAt"`
}
