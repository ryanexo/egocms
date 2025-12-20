package model

import (
    `dpcms/internal/infra/database`
)

type TokenBlacklist struct {
    database.Model
    UserId uint64 `gorm:"index;not null"`
    UUID   string `gorm:"column:uuid;unique;not null"`
}
