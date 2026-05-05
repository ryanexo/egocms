package domain

import (
    `time`
    
    `gorm.io/gorm`
)

type Category struct {
    id        uint64
    parentId  *uint64
    sequence  int64
    name      string
    path      string
    seo       CategorySEO
    createdAt time.Time
    updatedAt time.Time
    deletedAt gorm.DeletedAt
}

func NewCategory(
    name string,
    path string,
) Category {
    return Category{
        name:     name,
        path:     path,
        sequence: 0,
        seo:      NewCategorySeo(name, make([]string, 0), ""),
    }
}

func RestoreCategory(
    id uint64,
    name string,
    path string,
    seq int64,
    seo CategorySEO,
    createdAt time.Time,
    updatedAt time.Time,
    deletedAt gorm.DeletedAt,
) Category {
    return Category{
        id:        id,
        name:      name,
        path:      path,
        sequence:  seq,
        seo:       seo,
        createdAt: createdAt,
        updatedAt: updatedAt,
        deletedAt: deletedAt,
    }
}

func (s *Category) ID() uint64 {
    return s.id
}

func (s *Category) ParentID() *uint64 {
    return s.parentId
}

func (s *Category) Sequence() int64 {
    return s.sequence
}

func (s *Category) Name() string {
    return s.name
}

func (s *Category) Path() string {
    return s.path
}

func (s *Category) CreatedAt() time.Time {
    return s.createdAt
}

func (s *Category) UpdatedAt() time.Time {
    return s.updatedAt
}

func (s *Category) DeletedAt() gorm.DeletedAt {
    return s.deletedAt
}

func (s *Category) SEO() CategorySEO {
    return s.seo
}

func (s *Category) IsRoot() bool {
    return s.parentId == nil
}

func (s *Category) MoveTo(id uint64) {
    s.parentId = &id
}

func (s *Category) UpdateSEOMeta(seo CategorySEO) {
    s.seo = seo
}

func (s *Category) SetSequence(seq int64) {
    s.sequence = seq
}
