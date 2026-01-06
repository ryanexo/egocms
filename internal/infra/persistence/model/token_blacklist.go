package model

import (
    `time`
    
    `dpcms/internal/infra/persistence/datatype`
)

type TokenBlacklist struct {
    Base
    UserId  datatype.SafeUint64 `gorm:"index;not null"`
    UUID    string              `gorm:"column:uuid;unique;not null"`
    Expires time.Time           `gorm:"not null"`
}
