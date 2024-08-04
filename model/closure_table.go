package model

type ClosureTable interface {
    // INSERT INTO @@table (ancestor, descendant, distance)
    // SELECT @id, @id, 0
    // {{ if ancestor > 0 }}
    // UNION ALL
    // SELECT ancestor, @id, distance + 1
    // FROM @@table
    // WHERE descendant = @ancestor
    // {{ end }}
    CreateBranch(ancestor uint, id uint) error
    
    // DELETE FROM @@table
    // WHERE descendant = @id OR ancestor = @id
    RemoveBranch(id uint) error
}

type ClosureTableModel struct {
    Parent     uint `gorm:"-:migration;column:parent"`
    Ancestor   uint `gorm:"not null"`
    Descendant uint `gorm:"not null"`
    Distance   uint `gorm:"not null"`
}
