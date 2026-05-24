package model

import (
    `time`
)

type Article struct {
    Base
    
    ContentTypeID uint64
    AuthorID      uint64
    
    Url     *string
    Slug    *string
    Title   string
    Summary string
    Status  int8
}

type ArticleVersion struct {
    Base
    ArticleID uint64
    VersionNo uint64
    
    Title     string
    Content   string
    Summary   string
    ChangeLog string
    Published bool
}

type ArticlePublish struct {
    Base
    ArticleID uint64
    VersionID uint64
    
    PublishAt   time.Time
    PublishType int8
    
    Remark string
}

type ArticleCategoryRel struct {
    Base
    ArticleID  uint64
    CategoryID uint64
}

type ArticleTag struct {
    Base
    Name string
    Slug string
}

type ArticleTagRel struct {
    Base
    ArticleID uint64
    TagID     uint64
}

type ArticleComment struct {
    Base
    ParentID uint64
}

type ArticleStat struct {
    ArticleID    uint64
    ViewCount    int64
    LikeCount    int64
    CommentCount int64
    Score        float32
    UpdatedAt    time.Time
}
