package entity

type CategoryPage struct {
    ID         uint64 `gorm:"primaryKey"`
    CategoryID uint64 `gorm:"not null"`
    Title      string `gorm:"type:varchar(255);not null"`
    Data       string `gorm:"type:longtext;not null" json:"data,omitempty"`
}
