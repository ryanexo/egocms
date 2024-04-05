package entity

type Tag struct {
    ID     uint64 `gorm:"primaryKey"`
    Name   string `gorm:"type:varchar(32);not null"`
    Remark string `gorm:"type:varchar(255);not null"`
}
