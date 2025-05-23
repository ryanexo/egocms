package model

import (
    `gorm.io/gen`
)

type ClosureTable interface {
    // INSERT INTO @@table (ancestor, descendant, distance, parent)
    // SELECT @id, @id, 0, @parent
    // {{ if parent > 0 }}
    // UNION ALL
    // SELECT ancestor, @id, distance + 1, @parent
    // FROM @@table
    // WHERE descendant = @parent
    // {{ end }}
    CreateBranch(id int64, parent int64) error
    
    // DELETE FROM @@table WHERE ancestor IN (
    // SELECT descendant FROM @@table WHERE ancestor=@ancestor
    // ) OR descendant IN (
    // SELECT descendant FROM @@table WHERE ancestor=@ancestor
    // )
    RemoveBranch(ancestor int64) error
    
    // SELECT a.*, CASE WHEN b.ancestor IS NULL THEN 0 ELSE b.ancestor END AS parent FROM @@table AS a
    // LEFT JOIN @@table AS b ON a.descendant = b.descendant AND b.distance = 1
    // WHERE a.ancestor = @ancestor AND a.distance > 0
    // ORDER BY a.distance ASC
    FindDescendantByAncestor(ancestor int64) ([]*gen.T, error)
}

type ClosureTableModel struct {
    Parent     int64 `gorm:"not null"`
    Ancestor   int64 `gorm:"not null"`
    Descendant int64 `gorm:"not null"`
    Distance   int64 `gorm:"not null"`
}
