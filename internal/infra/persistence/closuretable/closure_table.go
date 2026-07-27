package closuretable

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "strings"
)

type ClosureTable struct {
    db  *sql.DB
    opt options
}

var (
    ErrInvalidNode  = errors.New("节点 ID 不合法")
    ErrCircularMove = errors.New("不能移动到自身或自己的子节点下")
)

var (
    ancestorColumn   = "ancestor"
    descendantColumn = "descendant"
    distanceColumn   = "distance"
)

func NewClosureTable(db *sql.DB, opts ...Option) (*ClosureTable, error) {
    opt := options{ModuleName: "default"}
    
    for _, option := range opts {
        option(&opt)
    }
    
    ct := &ClosureTable{db: db, opt: opt}
    err := ct.createTableIfNotExists()
    if err != nil {
        return nil, err
    }
    return ct, nil
}

func (c *ClosureTable) createTableIfNotExists() error {
    definitions := []string{
        fmt.Sprintf(`%s BIGINT NOT NULL`, ancestorColumn),
        fmt.Sprintf(`%s BIGINT NOT NULL`, descendantColumn),
        fmt.Sprintf(`%s BIGINT NOT NULL`, distanceColumn),
        fmt.Sprintf(`INDEX tree_rel(%s, %s)`, ancestorColumn, descendantColumn),
    }
    _, err := c.db.Exec(
        fmt.Sprintf(
            "CREATE TABLE%s IF NOT EXISTS (%s)",
            c.opt.tableName(),
            strings.Join(definitions, ", "),
        ),
    )
    return err
}

func (c *ClosureTable) CreateSubtree(ctx context.Context, id uint64, parentID uint64) error {
    if id == 0 {
        return ErrInvalidNode
    }
    
    tableName := c.opt.tableName()
    
    if parentID == 0 {
        _, err := c.db.ExecContext(
            ctx,
            fmt.Sprintf(
                "INSERT INTO %s (ancestor, descendant, distance) VALUES (?, ?, 0)",
                tableName,
            ),
            id,
            id,
        )
        return err
    }
    
    stmt := fmt.Sprintf(
        "INSERT INTO %s (ancestor, descendant, distance) "+
            "SELECT ancestor, ?, distance + 1 FROM %s WHERE descendant = ? "+
            "UNION ALL SELECT ?, ?, 0",
        tableName,
        tableName,
    )
    _, err := c.db.ExecContext(
        ctx,
        stmt,
        id,
        parentID,
        id,
        id,
    )
    return err
}

func (c *ClosureTable) Move(ctx context.Context, fromNode uint64, toNode uint64) error {
    if fromNode == 0 {
        return ErrInvalidNode
    }
    if fromNode == toNode {
        return ErrCircularMove
    }
    if toNode != 0 {
        isDescendant, err := c.exists(ctx, fromNode, toNode)
        if err != nil {
            return err
        }
        if isDescendant {
            return ErrCircularMove
        }
    }
    
    tx, err := c.db.BeginTx(ctx, nil)
    if err != nil {
        _ = tx.Rollback()
        return err
    }
    
    if err = c.unbindRelationships(ctx, fromNode); err != nil {
        return err
    }
    if toNode != 0 {
        if err = c.rebindRelationships(ctx, fromNode, toNode); err != nil {
            return err
        }
    }
    
    return tx.Commit()
}

func (c *ClosureTable) DropSubtree(ctx context.Context, root uint64) error {
    table := c.opt.tableName()
    _, err := c.db.ExecContext(
        ctx,
        fmt.Sprintf(
            "DELETE FROM %s WHERE %s IN ("+
                "SELECT d FROM (SELECT %s AS d FROM %s WHERE %s = ?) AS subtree_ancestors"+
                ") OR %s IN ("+
                "SELECT d FROM (SELECT %s AS d FROM %s WHERE %s = ?) AS subtree_descendants"+
                ")",
            table,
            ancestorColumn,
            descendantColumn,
            table,
            ancestorColumn,
            descendantColumn,
            descendantColumn,
            table,
            ancestorColumn,
        ),
        root,
        root,
    )
    return err
}

func (c *ClosureTable) exists(ctx context.Context, ancestor uint64, descendant uint64) (bool, error) {
    if ancestor == 0 || descendant == 0 {
        return false, nil
    }
    
    var exists uint8
    err := c.db.QueryRowContext(
        ctx,
        fmt.Sprintf(
            "SELECT 1 FROM %s WHERE ancestor = ? AND descendant = ? LIMIT 1",
            c.opt.tableName(),
        ),
        ancestor,
        descendant,
    ).Scan(&exists)
    if errors.Is(err, sql.ErrNoRows) {
        return false, nil
    }
    return err == nil, err
}

func (c *ClosureTable) unbindRelationships(ctx context.Context, root uint64) error {
    if root == 0 {
        return ErrInvalidNode
    }
    table := c.opt.tableName()
    _, err := c.db.ExecContext(
        ctx,
        fmt.Sprintf(
            "DELETE FROM %s WHERE %s IN ("+
                "SELECT d FROM (SELECT %s AS d FROM %s WHERE %s = ?) AS descendants"+
                ") AND %s IN ("+
                "SELECT a FROM (SELECT %s AS a FROM %s WHERE %s = ? AND %s <> ?) AS ancestors"+
                ")",
            table,
            descendantColumn,
            descendantColumn,
            table,
            ancestorColumn,
            ancestorColumn,
            ancestorColumn,
            table,
            descendantColumn,
            ancestorColumn,
        ),
        root,
        root,
        root,
    )
    
    return err
}

func (c *ClosureTable) rebindRelationships(ctx context.Context, root uint64, target uint64) error {
    table := c.opt.tableName()
    _, err := c.db.ExecContext(
        ctx,
        fmt.Sprintf(
            "INSERT INTO %s (%s, %s, %s) "+
                "SELECT A.%s, D.%s, A.%s + D.%s + 1 "+
                "FROM %s AS A CROSS JOIN %s AS D "+
                "WHERE A.%s = ? AND D.%s = ?",
            table,
            ancestorColumn,
            descendantColumn,
            distanceColumn,
            ancestorColumn,
            descendantColumn,
            distanceColumn,
            distanceColumn,
            table,
            table,
            descendantColumn,
            ancestorColumn,
        ),
        target,
        root,
    )
    return err
}
