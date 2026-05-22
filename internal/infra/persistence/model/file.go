package model

type File struct {
	Base
	UserID       uint64 `gorm:"index:idx_file_user_id;not null"`
	OriginalName string `gorm:"type:varchar(255);not null"`
	Ext          string `gorm:"type:varchar(32);not null"`
	Path         string `gorm:"type:varchar(255);index:idx_file_path;not null"`
	Size         int64  `gorm:"not null;index:idx_file_uniq;priority:2"`
	Driver       string `gorm:"type:varchar(32);index:idx_file_driver;default:local"`
	SHA256       string `gorm:"type:varchar(64);column:sha256;index:idx_file_uniq;priority:1;not null"`
	IsImage      bool   `gorm:"type:tinyint(1);default:0;index:idx_file_is_image;not null"`
}
