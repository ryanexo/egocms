package model

type Menu struct {
    Base
    ParentID uint64 `gorm:"not null" json:"parentId"`
    Name     string `gorm:"type:varchar(64);not null" json:"name"`
    Sequence int64  `gorm:"default:0;not null" json:"sequence"`
    Visible  bool   `gorm:"type:tinyint;default:1;not null" json:"visible"`
    URI      string `gorm:"type:varchar(255);not null" json:"URI"`
    Template string `gorm:"type:varchar(255);not null;" json:"template"`
    Remark   string `gorm:"type:varchar(255);default:'';not null" json:"remark"`
}

type MenuContext struct {
    ID uint64 `gorm:"primaryKey"`
    ClosureTableModel
}
