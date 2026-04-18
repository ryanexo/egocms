package domain

import `cms/internal/infra/persist/datatype`

type Category struct {
    id       datatype.SafeUint64
    parentId *datatype.SafeUint64
    sequence datatype.SafeInt64
    name     string
    path     string
    typ      int8
    url      *string
    visible  datatype.BoolInt8
    seo      *CategorySeo
    content  *string
}

func NewCategory(
    id datatype.SafeUint64,
    name string,
    path string,
) *Category {
    return &Category{
        id:       id,
        name:     name,
        path:     path,
        visible:  datatype.BoolInt8(1),
        typ:      0,
        sequence: 0,
        seo:      NewCategorySeo(name, make([]string, 0), ""),
    }
}

func (s *Category) IsRoot() bool {
    return s.parentId == nil
}

func (s *Category) IsCategory() bool {
    return s.typ == CategoryTypeNormal
}

func (s *Category) IsUrl() bool {
    return s.typ == CategoryTypeUrl
}

func (s *Category) IsPage() bool {
    return s.typ == CategoryTypePage
}

func (s *Category) SetParent() {}
