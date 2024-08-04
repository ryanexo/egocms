package model

import "time"

type UserLoginLog struct {
	ID        uint64    `gorm:"primaryKey"`
	UserID    uint64    `gorm:"not null;index"`
	CreatedIP string    `gorm:"type:varchar(39);not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
