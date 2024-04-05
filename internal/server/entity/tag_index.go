package entity

type TagIndex struct {
    ID        uint64 `gorm:"primaryKey"`
    TagID     uint64 `gorm:"uniqueIndex:relation;priority:1"`
    ContentID uint64 `gorm:"uniqueIndex:relation;priority:2"`
}
