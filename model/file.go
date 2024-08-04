package model

type File struct {
	ID       uint64 `gorm:"primaryKey"`
	Filename string `gorm:"type:varchar(255);"`
	Path     string `gorm:"type:varchar(255);not null"`
	Size     uint64 `gorm:"not null"`
	Ext      string `gorm:"type:varchar(255);not null"`
	IsImage  uint8  `gorm:"not null;index"`
	Md5      string `gorm:"type:varchar(32);not null;index:idx_hash"`
	Sha1     string `gorm:"type:varchar(40);not null;index:idx_hash"`
}
