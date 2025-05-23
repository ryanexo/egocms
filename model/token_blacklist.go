package model

import (
    `dpcms/packages/database`
)

type TokenBlacklist struct {
    database.Model
    UserId int64  `gorm:"index;not null"`
    UUID   string `gorm:"column:'uuid';unique;not null"`
}
