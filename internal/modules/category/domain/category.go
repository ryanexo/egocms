package domain

import (
    `regexp`
    
    `cms/internal/modules/category/internal/errno`
)

type Category struct {
    id       uint64
    parentID uint64
    sequence int64
    name     string
    path     string
    seo      CategorySEO
}

var pathRegexp = regexp.MustCompile(`^[a-zA-Z0-9]+([-_][a-zA-Z0-9]+)*[a-zA-Z0-9]$`)

func NewCategory(name string, path string) (Category, error) {
    if !pathRegexp.MatchString(path) {
        return Category{}, errno.ErrPathFormat
    }
    
    return Category{
        name: name,
        path: path,
    }, nil
}

func (c Category) ID() uint64 {
    return c.id
}

func (c Category) ParentID() uint64 {
    return c.parentID
}

func (c Category) Sequence() int64 {
    return c.sequence
}

func (c Category) Name() string {
    return c.name
}

func (c Category) Path() string {
    return c.path
}

func (c Category) SEO() CategorySEO {
    return c.seo
}

func (c Category) SetSEO(seo CategorySEO) {
    c.seo = seo
}

func (c Category) SetParent(v Category) {
    c.parentID = v.ID()
}

func (c Category) SetSequence(seq int64) {
    c.sequence = seq
}
