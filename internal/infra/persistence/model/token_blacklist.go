package model

import (
    `time`
)

type TokenBlacklist struct {
    Base
    UserID  uint64    `gorm:"index;not null"`
    UUID    string    `gorm:"column:uuid;unique;not null"`
    Expires time.Time `gorm:"not null"`
}
