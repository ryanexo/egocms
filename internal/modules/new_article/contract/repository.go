package contract

import (
    "context"
    "time"
    
    contentdomain "cms/internal/modules/contenttype/domain"
    "cms/internal/modules/new_article/domain"
)

type CreateArticle struct {
    AuthorID      uint64
    ContentTypeID *uint64
    URL           *string
    Slug          *string
}

type UpdateArticle struct {
    URL  *string
    Slug *string
}

type CreateRevision struct {
    ArticleID uint64
    Title     string
    Content   string
    Summary   string
    ChangeLog string
}

type RevisionDetail struct {
    ID        uint64
    VersionNo uint64
    Title     string
    Content   string
    Summary   string
    ChangeLog string
    CreatedAt time.Time
}

type ArticleDetail struct {
    ID                 uint64
    CreatedAt          time.Time
    UpdatedAt          time.Time
    ContentTypeID      *uint64
    AuthorID           uint64
    URL                *string
    Slug               *string
    Status             domain.Status
    CurrentVersionID   uint64
    PublishedVersionID *uint64
    CurrentRevision    RevisionDetail
    CategoryIDs        []uint64
    Tags               []string
    PublishedAt        *time.Time
    ContentSchemas     []contentdomain.EntrySchema
    ContentData        map[string]any
}

type PublishRecord struct {
    ArticleID uint64
    VersionID uint64
    PublishAt time.Time
    Direct    bool
}

type Repository interface {
    WithinTransaction(ctx context.Context, operation func(Repository) error) error
    CreateArticle(ctx context.Context, data CreateArticle) (uint64, error)
    CreateRevision(ctx context.Context, data CreateRevision) (uint64, error)
    SetCurrentRevision(ctx context.Context, articleID, revisionID uint64, title, summary string) error
    UpdateArticle(ctx context.Context, articleID uint64, data UpdateArticle) error
    ReplaceCategories(ctx context.Context, articleID uint64, categoryIDs []uint64) error
    ReplaceTags(ctx context.Context, articleID uint64, tags []string) error
    FindContentTypeID(ctx context.Context, articleID uint64) (*uint64, error)
    FindContentTypeSchemas(ctx context.Context, contentTypeID uint64) ([]contentdomain.EntrySchema, error)
    ReplaceContentData(ctx context.Context, articleID, contentTypeID uint64, entries contentdomain.EntrySet) error
    DeleteContentData(ctx context.Context, articleID uint64) error
    FindWorkflow(ctx context.Context, articleID uint64, forUpdate bool) (domain.Workflow, error)
    SaveWorkflow(ctx context.Context, workflow domain.Workflow) error
    CreatePublishRecord(ctx context.Context, data PublishRecord) error
    DeleteArticle(ctx context.Context, articleID uint64) error
    FindDetail(ctx context.Context, articleID uint64) (ArticleDetail, error)
}
