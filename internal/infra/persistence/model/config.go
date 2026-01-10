package model

type Config struct {
    Base
    Key   string `gorm:"type:varchar(64);not null;unique:uniq_cfg_key"`
    Value string `gorm:"type:varchar(255);not null"`
}
