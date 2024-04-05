package entity

type Group struct {
    ID         uint64 `gorm:"primaryKey"`
    Name       string `gorm:"type:varchar(64)"`
    Permission uint64 `gorm:"not null"`
}
