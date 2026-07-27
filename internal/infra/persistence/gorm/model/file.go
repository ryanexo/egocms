package model

type File struct {
	Base
	UserID       uint64 `gorm:"type:bigint unsigned"`
	OriginalName string `gorm:"type:varchar(255)"`
	Ext          string `gorm:"type:varchar(32)"`
	Path         string `gorm:"type:varchar(255)"`
	Size         int64  `gorm:"type:bigint"`
	Driver       *string `gorm:"type:varchar(32)"`
	SHA256       string `gorm:"type:varchar(64)"`
	IsImage      bool   `gorm:"type:tinyint(1)"`
}

func (File) TableName() string {
	return "file"
}
