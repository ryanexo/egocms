package model

import `dpcms/internal/infra/persistence/datatype`

type File struct {
    Base
    UserID       datatype.SafeUint64 `gorm:"index:idx_file_user_id;not null"`
    OriginalName string              `gorm:"type:varchar(255);not null"`
    Ext          string              `gorm:"type:varchar(32);not null"`
    Path         string              `gorm:"type:varchar(255);index:idx_file_path;not null"`
    Size         datatype.SafeInt64  `gorm:"not null;index:idx_file_uniq;priority:2"`
    Driver       string              `gorm:"type:varchar(32);index:idx_file_driver;default:local"`
    SHA256       string              `gorm:"type:varchar(64);index:idx_file_uniq;priority:1;not null"`
    IsImage      datatype.BoolInt8   `gorm:"type:tinyint(1);default:0;index:idx_file_is_image;not null"`
}
