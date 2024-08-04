package model

type Var struct {
	ID    uint64 `gorm:"primaryKey"`
	Name  string `gorm:"type:varchar(32);not null;uniqueIndex"`
	Value string `gorm:"type:varchar(255);not null;"`
}
