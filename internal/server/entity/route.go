package entity

type Route struct {
    ID  uint64 `gorm:"primaryKey"`
    URI string `gorm:"varchar(255)"`
}
