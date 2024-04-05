package entity

type AdditionFields struct {
    ID         uint64 `gorm:"primaryKey"`
    AdditionID uint64 `gorm:"not null;uniqueIndex:idx_addition;priority:1"`
    Name       string `gorm:"type:varchar(64);not null"`
    Identifier string `gorm:"type:varchar(64);not null;uniqueIndex:idx_addition;priority:2"`
    Type       uint8  `gorm:"not null"`
    Sequence   uint8  `gorm:"not null"`
    Required   uint8  `gorm:"not null"`
    Context    string `gorm:"type:text;not null"`
}
