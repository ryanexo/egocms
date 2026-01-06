package model

import "dpcms/internal/infra/persistence/datatype"

type Menu struct {
	Base
	ParentID datatype.SafeUint64 `gorm:"not null"`
	Type     int8                `gorm:"type:tinyint;not null"`
	Name     string              `gorm:"type:varchar(64);not null"`
	Sequence datatype.SafeInt64  `gorm:"default:0;not null"`
	Visible  datatype.BoolInt8   `gorm:"type:tinyint;default:1;not null"`
	URI      string              `gorm:"type:varchar(255);not null"`
	Resource string              `gorm:"type:varchar(255);not null"`
	Template string              `gorm:"type:varchar(255);not null;"`
	Remark   string              `gorm:"type:varchar(255);default:'';not null"`
	Action   []*MenuAction       `gorm:"foreignKey:MenuID;referenceKey:ID"`
}

type MenuContext struct {
	ID datatype.SafeUint64 `gorm:"primaryKey"`
	ClosureTableModel
}

type MenuAction struct {
	Base
	MenuID datatype.SafeUint64 `gorm:"not null"`
	Name   string              `gorm:"type:varchar(255);not null"`
	Action string              `gorm:"type:varchar(64);not null"`
}
