package entity

type ArticleData struct {
    ID        uint64 `gorm:"primaryKey"`
    ArticleID uint64 `gorm:"not null"`
    Title     string `gorm:"type:varchar(255);not null"`
    Data      string `gorm:"type:longtext;not null"`
}
