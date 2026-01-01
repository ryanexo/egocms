package model

import `dpcms/internal/infra/persistence/datatype`

type Menu struct {
    Base
    ParentID datatype.SafeUint64 `gorm:"not null" json:"parentId"`
    Name     string              `gorm:"type:varchar(64);not null" json:"name"`
    Sequence int64               `gorm:"default:0;not null" json:"sequence"`
    Visible  bool                `gorm:"type:tinyint;default:1;not null" json:"visible"`
    URI      string              `gorm:"type:varchar(255);not null" json:"uri"`
    Template string              `gorm:"type:varchar(255);not null;" json:"template"`
    Remark   string              `gorm:"type:varchar(255);default:'';not null" json:"remark"`
}

type MenuContext struct {
    ID datatype.SafeUint64 `gorm:"primaryKey"`
    ClosureTableModel
}
