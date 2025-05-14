package model

import (
    `gorm.io/gen`
)

type ClosureTable interface {
    // INSERT INTO @@table (ancestor, descendant, distance, parent)
    // SELECT @id, @id, 0
    // {{ if ancestor > 0 }}
    // UNION ALL
    // SELECT a.ancestor, @id, a.distance + 1
    // FROM @@table
    // WHERE descendant = @ancestor
    // {{ end }}
    CreateBranch(ancestor uint, id uint) error
    
    // DELETE FROM @@table
    // WHERE descendant = @id OR ancestor = @id
    RemoveBranch(id uint) error
    
    FindDescendantByAncestor(ancestor uint, id uint) ([]*gen.T, error)
}

type ClosureTableModel struct {
    Parent     uint `gorm:"not null"`
    Ancestor   uint `gorm:"not null"`
    Descendant uint `gorm:"not null"`
    Distance   uint `gorm:"not null"`
}
