package domain

import (
    `time`
    
    `cms/internal/modules/article/internal/errno`
    `cms/internal/pkg/stringx`
    
    `github.com/microcosm-cc/bluemonday`
)

var transitions = map[statusValue][]statusValue{
    StatusDraft:            {StatusPending},
    StatusPending:          {StatusPublished},
    StatusPublished:        {StatusOffline, StatusReject},
    StatusOffline:          {StatusPublished, StatusDraft},
    StatusReject:           {StatusPendingRepublish},
    StatusPendingRepublish: {StatusPublished},
}

type Article struct {
    id          uint64
    categoryID  uint64
    title       string
    description string
    content     string
    authorID    uint64
    flag        int16
    status      statusValue
    target      *string
    keywords    []string
    createdAt   time.Time
    updatedAt   time.Time
}

func NewArticle(title, content string, author uint64) Article {
    return Article{
        title:    title,
        content:  bluemonday.UGCPolicy().Sanitize(content),
        authorID: author,
        keywords: make([]string, 0),
    }
}

func (s *Article) CategoryID() uint64 {
    return s.categoryID
}

func (s *Article) ID() uint64 {
    return s.id
}

func (s *Article) Title() string {
    return s.title
}

func (s *Article) Description() string {
    if s.description == "" {
        return stringx.TruncateUTF8(s.content, 250)
    }
    return s.description
}

func (s *Article) Content() string {
    return s.content
}

func (s *Article) AuthorID() uint64 {
    return s.authorID
}

func (s *Article) Flag() int16 {
    return s.flag
}

func (s *Article) Status() int8 {
    return s.status.v
}

func (s *Article) Target() *string {
    return s.target
}

func (s *Article) Keywords() []string {
    return s.keywords
}

func (s *Article) CreatedAt() time.Time {
    return s.createdAt
}

func (s *Article) UpdatedAt() time.Time {
    return s.updatedAt
}

func (s *Article) toStatus(to statusValue) error {
    validStates, ok := transitions[s.status]
    if !ok {
        return errno.ErrInvalidStatusTransition
    }
    for _, validState := range validStates {
        if validState == to {
            s.status = to
            return nil
        }
    }
    return errno.ErrInvalidStatusTransition
}

func (s *Article) Submit() error {
    return s.toStatus(StatusPending)
}

func (s *Article) Publish() error {
    return s.toStatus(StatusPublished)
}

func (s *Article) Offline() error {
    return s.toStatus(StatusOffline)
}

func (s *Article) Reject() error {
    return s.toStatus(StatusOffline)
}

func (s *Article) Republish() error {
    return s.toStatus(StatusPendingRepublish)
}

func (s *Article) SetDescription(v string) {
    s.description = v
}

func (s *Article) SetFlag(v string) {}
