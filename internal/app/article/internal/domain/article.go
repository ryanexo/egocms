package domain

import "time"

type Article struct {
    id          uint64
    categoryID  uint64
    title       string
    description string
    authorID    uint64
    flag        int16
    status      int8
    target      *string
    keywords    []string
    createdAt   time.Time
    updatedAt   time.Time
}
