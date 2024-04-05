package entity

type Permission struct {
    ID          uint64 `gorm:"primaryKey"`
    Name        string `gorm:"type:varchar(255);not null"`
    Description string `gorm:"type:varchar(255);not null"`
    Value       uint64 `gorm:"not null"`
}
