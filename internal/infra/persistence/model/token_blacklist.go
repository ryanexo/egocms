package model

type TokenBlacklist struct {
    Base
    UserId uint64 `gorm:"index;not null"`
    UUID   string `gorm:"column:uuid;unique;not null"`
}
