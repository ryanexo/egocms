package model

type CategoryContext struct {
    ID         uint64
    Ancestor   uint64
    Descendant uint64
    Distance   uint64
}
