package model

import `dpcms/internal/infra/persistence/datatype`

type File struct {
    Base
    UserID  datatype.SafeUint64 `gorm:"index:idx_file_user_id;not null"`
    Name    string              `gorm:"type:varchar(255);not null"`
    Ext     string              `gorm:"type:varchar(32);not null"`
    Path    string              `gorm:"type:varchar(255);not null"`
    Size    datatype.SafeInt64  `gorm:"not null"`
    IsImage datatype.BoolInt8   `gorm:"type:tinyint(1);default:0"`
    SHA256  string              `gorm:"type:varchar(64);not null"`
    Driver  string              `gorm:"type:varchar(32);index:idx_file_driver;default:local"`
}
