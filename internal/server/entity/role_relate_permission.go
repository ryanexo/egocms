package entity

type RolePermission struct {
    ID         uint64 `gorm:"primaryKey"`
    RoleID     uint64 `gorm:"not null;index:idx_relation"`
    Permission uint64 `gorm:"not null"`
}
