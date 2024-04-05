package entity

type Role struct {
    ID          uint64 `gorm:"primaryKey"`
    Name        string `gorm:"type:varchar(255);not null"`
    Description string `gorm:"type:varchar(255);not null"`
}
