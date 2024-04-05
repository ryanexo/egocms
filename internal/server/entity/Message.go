package entity

import (
    `time`
)

type Message struct {
    ID                uint64    `gorm:"primaryKey"`
    Sender            uint64    `gorm:"not null"`
    Receiver          uint64    `gorm:"not null"`
    CreatedAt         time.Time `gorm:"autoCreateTime"`
    SenderDeletedAt   time.Time `gorm:"not null"`
    ReceiverDeletedAt time.Time `gorm:"not null"`
    ReadState         byte      `gorm:"not null"` // 0:已读, 1:未读
    DeleteState       byte      `gorm:"not null"` // 两位二进制表示 1为已删除
}
