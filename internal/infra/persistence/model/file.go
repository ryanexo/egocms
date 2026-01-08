package model

import `dpcms/internal/infra/persistence/datatype`

type File struct {
    Base
    Filename     string              `gorm:"type:varchar(255);not null"`
    Path         string              `gorm:"type:varchar(255);not null"`
    Size         datatype.SafeUint64 `gorm:"not null"`
    Driver       string              `gorm:"type:varchar(32);index:idx_file_driver;default:local"`
    Sha256       string              `gorm:"type:char(64);index:idx_file_hash;not null"`
    UploadRecord *UserUploadRecord   `gorm:"foreignKey:FileID;referenceKey:ID"`
}

type UserUploadRecord struct {
    Base
    UserID      datatype.SafeUint64 `gorm:"not null"`
    FileID      datatype.SafeUint64 `gorm:"not null"`
    Description string              `gorm:"type:varchar(255);default:''"`
}
