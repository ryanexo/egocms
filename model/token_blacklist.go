package model

type TokenBlacklist struct {
    Model
    UserId uint   `gorm:"index;not null"`
    UUID   string `gorm:"column:'uuid';unique;not null"`
}
