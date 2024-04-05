package entity

type Menu struct {
    ID       uint64 `gorm:"primaryKey"`
    RootID   uint64 `gorm:"not null"`
    ParentID uint64 `gorm:"not null;index"`
    RouteID  uint64 `gorm:"not null"`
    Sequence uint64 `gorm:"not null"`
    Name     string `gorm:"type:varchar(64);not null"`
    Display  byte   `gorm:"not null"`
}
