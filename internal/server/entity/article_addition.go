package entity

type ArticleAddition struct {
    ID             uint64 `gorm:"primaryKey"`
    ArticleID      uint64 `gorm:"not null"`
    AdditionID     uint64 `gorm:"not null"`
    AdditionDataID uint64 `gorm:"not null"`
}
