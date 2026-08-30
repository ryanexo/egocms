package new_article

import (
    "context"
    "strings"
    "time"
    "unicode/utf8"
    
    contentdomain "cms/internal/modules/contenttype/domain"
    "cms/internal/modules/new_article/contract"
    "cms/internal/modules/new_article/domain"
    `cms/internal/public/apitype`
    `cms/internal/public/jsontype`
)

const (
    maxTagCount  = 10
    maxTagLength = 64
)

type Service struct {
    repository contract.Repository
    now        func() time.Time
}

func NewArticleService(repository contract.Repository) *Service {
    return &Service{repository: repository, now: time.Now}
}

func (service *Service) CreateDraft(ctx context.Context, authorID uint64, params CreateParams) (jsontype.SafeUint64, error) {
    revision, err := domain.NewRevision(params.Title, params.Content, params.Summary, params.ChangeLog)
    if err != nil {
        return 0, mapDomainError(err)
    }
    contentTypeID := safeUint64Pointer(params.ContentTypeID)
    if contentTypeID != nil && *contentTypeID == 0 {
        return 0, ErrInvalidArticleRelation
    }
    categoryIDs, tags, err := normalizeRelations(params.CategoryIDs, params.Tags)
    if err != nil {
        return 0, err
    }
    
    var articleID uint64
    err = service.repository.WithinTransaction(ctx, func(repository contract.Repository) error {
        articleID, err = repository.CreateArticle(ctx, contract.CreateArticle{
            AuthorID:      authorID,
            ContentTypeID: contentTypeID,
            URL:           optionalString(params.URL),
            Slug:          optionalString(params.Slug),
        })
        if err != nil {
            return err
        }
        
        revisionID, createErr := repository.CreateRevision(ctx, contract.CreateRevision{
            ArticleID: articleID,
            Title:     revision.Title(),
            Content:   revision.Content(),
            Summary:   revision.Summary(),
            ChangeLog: revision.ChangeLog(),
        })
        if createErr != nil {
            return createErr
        }
        if createErr = repository.SetCurrentRevision(ctx, articleID, revisionID, revision.Title(), revision.Summary()); createErr != nil {
            return createErr
        }
        if createErr = repository.ReplaceCategories(ctx, articleID, categoryIDs); createErr != nil {
            return createErr
        }
        if createErr = repository.ReplaceTags(ctx, articleID, tags); createErr != nil {
            return createErr
        }
        return service.replaceContentData(ctx, repository, articleID, contentTypeID, params.ContentData)
    })
    if err != nil {
        return 0, mapDomainError(err)
    }
    return jsontype.SafeUint64(articleID), nil
}

func (service *Service) UpdateDraft(ctx context.Context, params UpdateParams) error {
    revision, err := domain.NewRevision(params.Title, params.Content, params.Summary, params.ChangeLog)
    if err != nil {
        return mapDomainError(err)
    }
    categoryIDs, tags, err := normalizeRelations(params.CategoryIDs, params.Tags)
    if err != nil {
        return err
    }
    
    articleID := params.ID.Uint64()
    err = service.repository.WithinTransaction(ctx, func(repository contract.Repository) error {
        workflow, findErr := repository.FindWorkflow(ctx, articleID, true)
        if findErr != nil {
            return findErr
        }
        if editErr := workflow.CanEdit(); editErr != nil {
            return editErr
        }
        contentTypeID, findErr := repository.FindContentTypeID(ctx, articleID)
        if findErr != nil {
            return findErr
        }
        
        revisionID, createErr := repository.CreateRevision(ctx, contract.CreateRevision{
            ArticleID: articleID,
            Title:     revision.Title(),
            Content:   revision.Content(),
            Summary:   revision.Summary(),
            ChangeLog: revision.ChangeLog(),
        })
        if createErr != nil {
            return createErr
        }
        if createErr = workflow.SetCurrentVersion(revisionID); createErr != nil {
            return createErr
        }
        if createErr = repository.UpdateArticle(ctx, articleID, contract.UpdateArticle{
            URL:  optionalString(params.URL),
            Slug: optionalString(params.Slug),
        }); createErr != nil {
            return createErr
        }
        if createErr = repository.SetCurrentRevision(ctx, articleID, revisionID, revision.Title(), revision.Summary()); createErr != nil {
            return createErr
        }
        if createErr = repository.ReplaceCategories(ctx, articleID, categoryIDs); createErr != nil {
            return createErr
        }
        if createErr = repository.ReplaceTags(ctx, articleID, tags); createErr != nil {
            return createErr
        }
        return service.replaceContentData(ctx, repository, articleID, contentTypeID, params.ContentData)
    })
    return mapDomainError(err)
}

func (service *Service) Submit(ctx context.Context, articleID jsontype.SafeUint64, canPublishDirect bool) error {
    return service.changeWorkflow(ctx, articleID.Uint64(), canPublishDirect, func(workflow *domain.Workflow) error {
        return workflow.Submit(canPublishDirect)
    })
}

func (service *Service) Publish(ctx context.Context, articleID jsontype.SafeUint64) error {
    return service.changeWorkflow(ctx, articleID.Uint64(), false, func(workflow *domain.Workflow) error {
        return workflow.Publish()
    })
}

func (service *Service) Reject(ctx context.Context, articleID jsontype.SafeUint64) error {
    return service.changeWorkflow(ctx, articleID.Uint64(), false, func(workflow *domain.Workflow) error {
        return workflow.Reject()
    })
}

func (service *Service) Offline(ctx context.Context, articleID jsontype.SafeUint64) error {
    return service.changeWorkflow(ctx, articleID.Uint64(), false, func(workflow *domain.Workflow) error {
        return workflow.Offline()
    })
}

func (service *Service) Republish(ctx context.Context, articleID jsontype.SafeUint64, canPublishDirect bool) error {
    return service.changeWorkflow(ctx, articleID.Uint64(), canPublishDirect, func(workflow *domain.Workflow) error {
        return workflow.Republish(canPublishDirect)
    })
}

func (service *Service) changeWorkflow(
    ctx context.Context,
    articleID uint64,
    direct bool,
    change func(*domain.Workflow) error,
) error {
    err := service.repository.WithinTransaction(ctx, func(repository contract.Repository) error {
        workflow, findErr := repository.FindWorkflow(ctx, articleID, true)
        if findErr != nil {
            return findErr
        }
        previousStatus := workflow.Status()
        previousPublishedVersionID := workflow.PublishedVersionID()
        if changeErr := change(&workflow); changeErr != nil {
            return changeErr
        }
        if saveErr := repository.SaveWorkflow(ctx, workflow); saveErr != nil {
            return saveErr
        }
        
        if shouldRecordPublish(previousStatus, previousPublishedVersionID, workflow) {
            return repository.CreatePublishRecord(ctx, contract.PublishRecord{
                ArticleID: workflow.ArticleID(),
                VersionID: workflow.CurrentVersionID(),
                PublishAt: service.now(),
                Direct:    direct,
            })
        }
        return nil
    })
    return mapDomainError(err)
}

func (service *Service) Delete(ctx context.Context, articleID jsontype.SafeUint64) error {
    err := service.repository.WithinTransaction(ctx, func(repository contract.Repository) error {
        if _, findErr := repository.FindWorkflow(ctx, articleID.Uint64(), true); findErr != nil {
            return findErr
        }
        return repository.DeleteArticle(ctx, articleID.Uint64())
    })
    return mapDomainError(err)
}

func (service *Service) FindByID(ctx context.Context, articleID jsontype.SafeUint64) (*Article, error) {
    detail, err := service.repository.FindDetail(ctx, articleID.Uint64())
    if err != nil {
        return nil, mapDomainError(err)
    }
    
    categoryIDs := make([]jsontype.SafeUint64, 0, len(detail.CategoryIDs))
    for _, categoryID := range detail.CategoryIDs {
        categoryIDs = append(categoryIDs, jsontype.SafeUint64(categoryID))
    }
    contentSchemas := make([]ContentSchema, 0, len(detail.ContentSchemas))
    for _, schema := range detail.ContentSchemas {
        contentSchemas = append(contentSchemas, ContentSchema{
            FieldKey:    schema.FieldKey,
            FieldName:   schema.FieldName,
            Description: schema.Description,
            Type:        schema.Type,
            Sequence:    schema.Sequence,
            Required:    schema.Required,
            Visible:     schema.Visible,
            Enabled:     schema.Enabled,
        })
    }
    return &Article{
        Base: apitype.Base{
            ID:        jsontype.SafeUint64(detail.ID),
            CreatedAt: detail.CreatedAt,
            UpdatedAt: detail.UpdatedAt,
        },
        ContentTypeID:      datatypePointer(detail.ContentTypeID),
        AuthorID:           jsontype.SafeUint64(detail.AuthorID),
        URL:                detail.URL,
        Slug:               detail.Slug,
        Status:             int8(detail.Status),
        CurrentVersionID:   jsontype.SafeUint64(detail.CurrentVersionID),
        PublishedVersionID: datatypePointer(detail.PublishedVersionID),
        CurrentRevision: Revision{
            ID:        jsontype.SafeUint64(detail.CurrentRevision.ID),
            VersionNo: jsontype.SafeUint64(detail.CurrentRevision.VersionNo),
            Title:     detail.CurrentRevision.Title,
            Content:   detail.CurrentRevision.Content,
            Summary:   detail.CurrentRevision.Summary,
            ChangeLog: detail.CurrentRevision.ChangeLog,
            CreatedAt: detail.CurrentRevision.CreatedAt,
        },
        CategoryIDs:    categoryIDs,
        Tags:           detail.Tags,
        PublishedAt:    detail.PublishedAt,
        ContentSchemas: contentSchemas,
        ContentData:    detail.ContentData,
    }, nil
}

func (service *Service) replaceContentData(
    ctx context.Context,
    repository contract.Repository,
    articleID uint64,
    contentTypeID *uint64,
    input map[string]any,
) error {
    if contentTypeID == nil {
        if len(input) > 0 {
            return ErrContentTypeRequired
        }
        return repository.DeleteContentData(ctx, articleID)
    }
    schemas, err := repository.FindContentTypeSchemas(ctx, *contentTypeID)
    if err != nil {
        return err
    }
    entries, err := contentdomain.BuildEntrySet(schemas, input)
    if err != nil {
        return ErrInvalidContentTypeData.Format(err.Error())
    }
    return repository.ReplaceContentData(ctx, articleID, *contentTypeID, entries)
}

func shouldRecordPublish(previousStatus domain.Status, previousVersionID *uint64, workflow domain.Workflow) bool {
    if workflow.Status() != domain.StatusPublished {
        return false
    }
    currentPublishedVersionID := workflow.PublishedVersionID()
    if currentPublishedVersionID == nil {
        return false
    }
    return previousStatus != domain.StatusPublished || !sameUint64Pointer(previousVersionID, currentPublishedVersionID)
}

func normalizeRelations(categoryValues []jsontype.SafeUint64, tagValues []string) ([]uint64, []string, error) {
    categoryIDs := make([]uint64, 0, len(categoryValues))
    categorySet := make(map[uint64]struct{}, len(categoryValues))
    for _, value := range categoryValues {
        categoryID := value.Uint64()
        if categoryID == 0 {
            return nil, nil, ErrInvalidArticleRelation
        }
        if _, exists := categorySet[categoryID]; exists {
            continue
        }
        categorySet[categoryID] = struct{}{}
        categoryIDs = append(categoryIDs, categoryID)
    }
    
    if len(tagValues) > maxTagCount {
        return nil, nil, ErrInvalidArticleRelation
    }
    tags := make([]string, 0, len(tagValues))
    tagSet := make(map[string]struct{}, len(tagValues))
    for _, value := range tagValues {
        tag := strings.TrimSpace(value)
        if tag == "" || utf8.RuneCountInString(tag) > maxTagLength {
            return nil, nil, ErrInvalidArticleRelation
        }
        key := strings.ToLower(tag)
        if _, exists := tagSet[key]; exists {
            continue
        }
        tagSet[key] = struct{}{}
        tags = append(tags, tag)
    }
    return categoryIDs, tags, nil
}

func optionalString(value string) *string {
    value = strings.TrimSpace(value)
    if value == "" {
        return nil
    }
    return &value
}

func safeUint64Pointer(value *jsontype.SafeUint64) *uint64 {
    if value == nil {
        return nil
    }
    result := value.Uint64()
    return &result
}

func datatypePointer(value *uint64) *jsontype.SafeUint64 {
    if value == nil {
        return nil
    }
    result := jsontype.SafeUint64(*value)
    return &result
}

func sameUint64Pointer(left, right *uint64) bool {
    if left == nil || right == nil {
        return left == nil && right == nil
    }
    return *left == *right
}
