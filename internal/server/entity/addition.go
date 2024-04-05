package entity

type Addition struct {
    ID           uint64 `gorm:"primaryKey"`
    Name         string `gorm:"type:varchar(255);not null"`
    ContentTable string `gorm:"type:varchar(64);not null;uniqueIndex"`
}
