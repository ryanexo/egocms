package model

import `cms/internal/infra/persistence/datatype`

type ClosureTable interface {
    // INSERT INTO @@table (ancestor, descendant, distance)
    // SELECT @id, @id, 0, @ancestor
    // {{ if ancestor > 0 }}
    // UNION ALL
    // SELECT ancestor, @id, distance + 1, @ancestor
    // FROM @@table
    // WHERE descendant = @ancestor
    // {{ end }}
    CreateSubtree(id uint64, ancestor uint64) error
    
    // DELETE FROM @@table WHERE ancestor IN (
    // SELECT descendant FROM @@table WHERE ancestor=@ancestor
    // ) OR descendant IN (
    // SELECT descendant FROM @@table WHERE ancestor=@ancestor
    // )
    DropSubtree(ancestor uint64) error
    
    // DELETE FROM @@table WHERE descendant IN (
    // SELECT d FROM ( SELECT descendant AS d FROM @@table WHERE ancestor=@ancestor ) AS DCT
    // ) AND ancestor IN (
    // SELECT a FROM ( SELECT ancestor AS a FROM @@table WHERE descendant=@ancestor AND ancestor<>@ancestor ) AS ACT
    // )
    UnbindRelationships(ancestor uint64) error
    
    // INSERT INTO @@table(ancestor, descendant, distance) SELECT A.ancestor, D.descendant, A.distance + D.distance + 1 FROM @@table AS A
    // CROSS JOIN @@table AS D WHERE A.descendant = @target AND D.ancestor = @ancestor
    ReBindRelationships(ancestor uint64, target uint64) error
}

type ClosureTableModel struct {
    Ancestor   datatype.SafeUint64 `gorm:"not null;index:,composite:ct;priority:1"`
    Descendant datatype.SafeUint64 `gorm:"not null;index:,composite:ct;priority:2"`
    Distance   datatype.SafeUint64 `gorm:"not null"`
}
