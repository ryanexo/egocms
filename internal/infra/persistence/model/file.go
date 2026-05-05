package model

type File struct {
    Base
    UserID       uint64 `gorm:"index:idx_file_user_id;not null"`
    OriginalName string `gorm:"type:varchar(255);not null;comment:'原始文件名'"`
    Ext          string `gorm:"type:varchar(32);not null;comment:'扩展名'"`
    Path         string `gorm:"type:varchar(255);index:idx_file_path;not null;comment:'文件路径'"`
    Size         int64  `gorm:"not null;index:idx_file_uniq;priority:2;comment:'文件尺寸'"`
    Driver       string `gorm:"type:varchar(32);index:idx_file_driver;default:local;comment:'文件驱动'"`
    SHA256       string `gorm:"type:varchar(64);index:idx_file_uniq;priority:1;not null"`
    IsImage      int8   `gorm:"type:tinyint(1);default:0;index:idx_file_is_image;not null"`
}
