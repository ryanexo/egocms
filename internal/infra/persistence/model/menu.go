package model

type Menu struct {
	Base
	ParentID    uint64 `gorm:"not null"`
	Type        int8   `gorm:"type:tinyint;not null"`
	Name        string `gorm:"type:varchar(64);not null"`
	Affix       bool   `gorm:"type:tinyint;default:0;not null"`
	Icon        string `gorm:"type:varchar(64)"`
	ExternalURL string `gorm:"type:varchar(255)"`
	Sequence    int64  `gorm:"default:0;not null"`
	Visible     bool   `gorm:"type:tinyint;default:1;not null"`
	URI         string `gorm:"type:varchar(255);not null"`
	Template    string `gorm:"type:varchar(255);not null;"`
	Remark      string `gorm:"type:varchar(255);default:'';not null"`
}

type MenuContext struct {
	ID uint64 `gorm:"primaryKey"`
	ClosureTableModel
}
