package model

import (
	"time"
)

type PointsLog struct {
	ID      uint64    `gorm:"primaryKey"`
	UserID  uint64    `gorm:"not null;index:idx_user;priority:1"`
	AdminID uint64    `gorm:"not null;index:idx_user;priority:2"`
	Type    byte      `gorm:"not null"` // 0:增加, 1:减少
	Time    time.Time `gorm:"autoCreateTime"`
	Count   uint64    `gorm:"not null"`
	Remark  string    `gorm:"type:varchar(255);not null;"`
}
