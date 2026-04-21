package domain

import (
    `cms/internal/infra/persist/datatype`
    valueobject2 `cms/internal/kernel/valueobject`
)

type Category struct {
    id       uint64
    parentId *uint64
    sequence datatype.SafeInt64
    name     string
    path     valueobject2.URLPath
    typ      int8
    url      *valueobject2.URL
    visible  int8
    seo      CategorySeo
    content  *valueobject2.HTMLContent
}

func NewCategory(
    name string,
    path valueobject2.URLPath,
) *Category {
    return &Category{
        name:     name,
        path:     path,
        visible:  1,
        typ:      0,
        sequence: 0,
        seo:      NewCategorySeo(name, make([]string, 0), ""),
    }
}

func (s *Category) Id() uint64 {
    return s.id
}

func (s *Category) ParentId() *uint64 {
    return s.parentId
}

func (s *Category) Sequence() datatype.SafeInt64 {
    return s.sequence
}

func (s *Category) Name() string {
    return s.name
}

func (s *Category) Path() valueobject2.URLPath {
    return s.path
}

func (s *Category) Typ() int8 {
    return s.typ
}

func (s *Category) URL() (*valueobject2.URL, bool) {
    return s.url, s.typ == TypeURL
}

func (s *Category) Visible() int8 {
    return s.visible
}

func (s *Category) Seo() CategorySeo {
    return s.seo
}

func (s *Category) Content() (*valueobject2.HTMLContent, bool) {
    return s.content, s.typ == TypePage
}

func (s *Category) IsRoot() bool {
    return s.parentId == nil
}

func (s *Category) MoveTo(id uint64) {
    s.parentId = &id
}

func (s *Category) Show() {
    s.visible = 1
}

func (s *Category) Hide() {
    s.visible = 0
}

func (s *Category) UpdateSEOMeta(seo CategorySeo) {
    s.seo = seo
}

func (s *Category) AsPage(content valueobject2.HTMLContent) {
    s.typ = TypePage
    s.content = &content
    s.url = nil
}

func (s *Category) AsURL(url valueobject2.URL) {
    s.typ = TypeURL
    s.url = &url
    s.content = nil
}

func (s *Category) AsCategory() {
    s.typ = TypeCategory
    s.url = nil
    s.content = nil
}

func (s *Category) SetSequence(seq datatype.SafeInt64) {
    s.sequence = seq
}
