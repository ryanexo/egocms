package entity

import (
    `gorm.io/gorm`
)

type FileRecord struct {
    gorm.Model
    UserID        uint64 `gorm:"not null;index:idx_record;priority:1"`
    FileID        uint64 `gorm:"not null;index:idx_record;priority:2"`
    Remark        string `gorm:"type:varchar(255);not null"`
    DownloadCount uint64 `gorm:"not null"`
    Context       string `gorm:"not null"` // 备用字段，JSON格式
}
