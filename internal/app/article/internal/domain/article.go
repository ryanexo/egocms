package domain

import "time"

type Article struct {
	ID        uint64
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}