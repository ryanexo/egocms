package new_article

import (
    "context"
    "encoding/json"
    "errors"
    "time"
    
    contentdomain "cms/internal/modules/contenttype/domain"
    "cms/internal/modules/new_article/contract"
    "cms/internal/modules/new_article/domain"
    `cms/internal/public/model`
    
    "gorm.io/datatypes"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"
)

type repository struct {
    db *gorm.DB
}

var _ contract.Repository = (*repository)(nil)

func NewRepository(db *gorm.DB) contract.Repository {
    return &repository{db: db}
}

func (repo *repository) WithinTransaction(ctx context.Context, operation func(contract.Repository) error) error {
    return repo.db.WithContext(ctx).Transaction(func(transaction *gorm.DB) error {
        return operation(&repository{db: transaction})
    })
}

func (repo *repository) CreateArticle(ctx context.Context, data contract.CreateArticle) (uint64, error) {
    record := articleRecord{
        ContentTypeID: data.ContentTypeID,
        AuthorID:      data.AuthorID,
        URL:           data.URL,
        Slug:          data.Slug,
        Status:        int8(domain.StatusDraft),
    }
    if err := repo.db.WithContext(ctx).Create(&record).Error; err != nil {
        return 0, err
    }
    return record.ID, nil
}

func (repo *repository) CreateRevision(ctx context.Context, data contract.CreateRevision) (uint64, error) {
    var versionNo uint64
    err := repo.db.WithContext(ctx).
        Model(&articleVersionRecord{}).
        Where("article_id = ?", data.ArticleID).
        Select("COALESCE(MAX(version_no), 0) + 1").
        Scan(&versionNo).Error
    if err != nil {
        return 0, err
    }
    
    record := articleVersionRecord{
        ArticleID: data.ArticleID,
        VersionNo: versionNo,
        Title:     data.Title,
        Content:   data.Content,
        Summary:   data.Summary,
        ChangeLog: data.ChangeLog,
    }
    if err := repo.db.WithContext(ctx).Create(&record).Error; err != nil {
        return 0, err
    }
    return record.ID, nil
}

func (repo *repository) SetCurrentRevision(ctx context.Context, articleID, revisionID uint64, title, summary string) error {
    return repo.db.WithContext(ctx).
        Model(&articleRecord{}).
        Where("id = ?", articleID).
        Updates(map[string]any{
            "current_version_id": revisionID,
            "title":              title,
            "summary":            summary,
        }).Error
}

func (repo *repository) UpdateArticle(ctx context.Context, articleID uint64, data contract.UpdateArticle) error {
    return repo.db.WithContext(ctx).
        Model(&articleRecord{}).
        Where("id = ?", articleID).
        Updates(map[string]any{
            "url":  data.URL,
            "slug": data.Slug,
        }).Error
}

func (repo *repository) ReplaceCategories(ctx context.Context, articleID uint64, categoryIDs []uint64) error {
    database := repo.db.WithContext(ctx)
    if err := database.Where("article_id = ?", articleID).Delete(&articleCategoryRecord{}).Error; err != nil {
        return err
    }
    if len(categoryIDs) == 0 {
        return nil
    }
    
    records := make([]articleCategoryRecord, 0, len(categoryIDs))
    for _, categoryID := range categoryIDs {
        records = append(records, articleCategoryRecord{ArticleID: articleID, CategoryID: categoryID})
    }
    return database.Create(&records).Error
}

func (repo *repository) ReplaceTags(ctx context.Context, articleID uint64, tags []string) error {
    database := repo.db.WithContext(ctx)
    if err := database.Where("article_id = ?", articleID).Delete(&articleTagRelationRecord{}).Error; err != nil {
        return err
    }
    if len(tags) == 0 {
        return nil
    }
    
    relations := make([]articleTagRelationRecord, 0, len(tags))
    for _, tag := range tags {
        tagRecord, err := repo.findOrCreateTag(ctx, tag)
        if err != nil {
            return err
        }
        relations = append(relations, articleTagRelationRecord{ArticleID: articleID, TagID: tagRecord.ID})
    }
    return database.Create(&relations).Error
}

func (repo *repository) FindContentTypeID(ctx context.Context, articleID uint64) (*uint64, error) {
    var record articleRecord
    if err := repo.db.WithContext(ctx).
        Select("content_type_id").
        Where("id = ?", articleID).
        First(&record).Error; err != nil {
        return nil, err
    }
    return cloneUint64(record.ContentTypeID), nil
}

func (repo *repository) FindContentTypeSchemas(ctx context.Context, contentTypeID uint64) ([]contentdomain.EntrySchema, error) {
    database := repo.db.WithContext(ctx)
    if err := database.Where("id = ?", contentTypeID).First(&model.ContentType{}).Error; err != nil {
        return nil, err
    }
    
    var records []model.ContentTypeSchema
    if err := database.
        Where("content_type_id = ?", contentTypeID).
        Order("sequence ASC, id ASC").
        Find(&records).Error; err != nil {
        return nil, err
    }
    
    schemas := make([]contentdomain.EntrySchema, 0, len(records))
    for _, record := range records {
        enumOptions, err := contentdomain.DecodeEnumOptions(record.EnumOptions)
        if err != nil {
            return nil, err
        }
        schemas = append(schemas, contentdomain.EntrySchema{
            FieldKey:    record.FieldKey,
            FieldName:   record.FieldName,
            Description: record.Description,
            Type:        record.Type,
            Sequence:    record.Sequence,
            MinLen:      record.MinLen,
            MaxLen:      record.MaxLen,
            MinValue:    record.MinValue,
            MaxValue:    record.MaxValue,
            MinTime:     record.MinTime,
            MaxTime:     record.MaxTime,
            Pattern:     record.Pattern,
            EnumOptions: enumOptions,
            Required:    boolValue(record.Required, false),
            Visible:     boolValue(record.Visible, true),
            Enabled:     boolValue(record.Enable, true),
        })
    }
    return schemas, nil
}

func (repo *repository) ReplaceContentData(
    ctx context.Context,
    articleID, contentTypeID uint64,
    entries contentdomain.EntrySet,
) error {
    database := repo.db.WithContext(ctx)
    if err := repo.DeleteContentData(ctx, articleID); err != nil {
        return err
    }
    
    encodedData, err := json.Marshal(entries.Data)
    if err != nil {
        return err
    }
    if err := database.Create(&model.ContentTypeEntries{
        ArticleID:     articleID,
        ContentTypeID: contentTypeID,
        Data:          datatypes.JSON(encodedData),
    }).Error; err != nil {
        return err
    }
    
    if len(entries.Values) == 0 {
        return nil
    }
    values := make([]model.ContentFieldValues, 0, len(entries.Values))
    for _, entry := range entries.Values {
        values = append(values, model.ContentFieldValues{
            ModelID:     contentTypeID,
            ArticleID:   articleID,
            FieldKey:    entry.FieldKey,
            Type:        entry.Type,
            StringValue: entry.StringValue,
            NumberValue: entry.NumberValue,
            BoolValue:   entry.BoolValue,
            TimeValue:   entry.TimeValue,
        })
    }
    return database.Create(&values).Error
}

func (repo *repository) DeleteContentData(ctx context.Context, articleID uint64) error {
    database := repo.db.WithContext(ctx).Unscoped()
    if err := database.Where("article_id = ?", articleID).Delete(&model.ContentFieldValues{}).Error; err != nil {
        return err
    }
    return database.Where("article_id = ?", articleID).Delete(&model.ContentTypeEntries{}).Error
}

func (repo *repository) findOrCreateTag(ctx context.Context, name string) (articleTagRecord, error) {
    database := repo.db.WithContext(ctx)
    var record articleTagRecord
    err := database.Unscoped().Where("name = ?", name).First(&record).Error
    switch {
    case err == nil:
        if record.DeletedAt.Valid {
            record.DeletedAt = gorm.DeletedAt{}
            err = database.Unscoped().Model(&record).Update("deleted_at", nil).Error
        }
        return record, err
    case !errors.Is(err, gorm.ErrRecordNotFound):
        return articleTagRecord{}, err
    }
    
    record = articleTagRecord{Name: name}
    if err := database.Clauses(clause.OnConflict{DoNothing: true}).Create(&record).Error; err != nil {
        return articleTagRecord{}, err
    }
    if record.ID == 0 {
        if err := database.Unscoped().Where("name = ?", name).First(&record).Error; err != nil {
            return articleTagRecord{}, err
        }
    }
    return record, nil
}

func (repo *repository) FindWorkflow(ctx context.Context, articleID uint64, forUpdate bool) (domain.Workflow, error) {
    database := repo.db.WithContext(ctx).
        Select("id", "status", "current_version_id", "published_version_id").
        Where("id = ?", articleID)
    if forUpdate {
        database = database.Clauses(clause.Locking{Strength: "UPDATE"})
    }
    
    var record articleRecord
    if err := database.First(&record).Error; err != nil {
        return domain.Workflow{}, err
    }
    currentVersionID := uint64(0)
    if record.CurrentVersionID != nil {
        currentVersionID = *record.CurrentVersionID
    }
    return domain.RestoreWorkflow(
        record.ID,
        domain.Status(record.Status),
        currentVersionID,
        record.PublishedVersionID,
    )
}

func (repo *repository) SaveWorkflow(ctx context.Context, workflow domain.Workflow) error {
    return repo.db.WithContext(ctx).
        Model(&articleRecord{}).
        Where("id = ?", workflow.ArticleID()).
        Updates(map[string]any{
            "status":               int8(workflow.Status()),
            "published_version_id": workflow.PublishedVersionID(),
        }).Error
}

func (repo *repository) CreatePublishRecord(ctx context.Context, data contract.PublishRecord) error {
    publishType := int8(0)
    if data.Direct {
        publishType = 1
    }
    return repo.db.WithContext(ctx).Create(&articlePublishRecord{
        ArticleID:   data.ArticleID,
        VersionID:   data.VersionID,
        PublishAt:   data.PublishAt,
        PublishType: publishType,
    }).Error
}

func (repo *repository) DeleteArticle(ctx context.Context, articleID uint64) error {
    database := repo.db.WithContext(ctx)
    if err := database.Where("id = ?", articleID).Delete(&articleRecord{}).Error; err != nil {
        return err
    }
    if err := database.Where("article_id = ?", articleID).Delete(&articleVersionRecord{}).Error; err != nil {
        return err
    }
    if err := database.Where("article_id = ?", articleID).Delete(&articlePublishRecord{}).Error; err != nil {
        return err
    }
    if err := database.Where("article_id = ?", articleID).Delete(&articleCategoryRecord{}).Error; err != nil {
        return err
    }
    if err := database.Where("article_id = ?", articleID).Delete(&articleTagRelationRecord{}).Error; err != nil {
        return err
    }
    return repo.DeleteContentData(ctx, articleID)
}

func (repo *repository) FindDetail(ctx context.Context, articleID uint64) (contract.ArticleDetail, error) {
    database := repo.db.WithContext(ctx)
    var article articleRecord
    if err := database.Where("id = ?", articleID).First(&article).Error; err != nil {
        return contract.ArticleDetail{}, err
    }
    if article.CurrentVersionID == nil {
        return contract.ArticleDetail{}, domain.ErrMissingCurrentVersion
    }
    
    var revision articleVersionRecord
    if err := database.
        Where("id = ? AND article_id = ?", *article.CurrentVersionID, article.ID).
        First(&revision).Error; err != nil {
        return contract.ArticleDetail{}, err
    }
    
    categoryIDs := make([]uint64, 0)
    if err := database.Model(&articleCategoryRecord{}).
        Where("article_id = ?", articleID).
        Order("category_id ASC").
        Pluck("category_id", &categoryIDs).Error; err != nil {
        return contract.ArticleDetail{}, err
    }
    
    tags := make([]string, 0)
    if err := database.Table("article_tag").
        Select("article_tag.name").
        Joins("JOIN article_tag_rel ON article_tag_rel.tag_id = article_tag.id").
        Where("article_tag_rel.article_id = ? AND article_tag.deleted_at IS NULL", articleID).
        Order("article_tag.name ASC").
        Scan(&tags).Error; err != nil {
        return contract.ArticleDetail{}, err
    }
    
    var publishedAt *time.Time
    var publish articlePublishRecord
    err := database.Where("article_id = ?", articleID).
        Order("publish_at DESC, id DESC").
        First(&publish).Error
    if err == nil {
        publishedAt = &publish.PublishAt
    } else if !errors.Is(err, gorm.ErrRecordNotFound) {
        return contract.ArticleDetail{}, err
    }
    
    contentSchemas := make([]contentdomain.EntrySchema, 0)
    contentData := make(map[string]any)
    if article.ContentTypeID != nil {
        contentSchemas, err = repo.FindContentTypeSchemas(ctx, *article.ContentTypeID)
        if err != nil {
            return contract.ArticleDetail{}, err
        }
        var entries model.ContentTypeEntries
        err = database.Where("article_id = ? AND content_type_id = ?", articleID, *article.ContentTypeID).
            First(&entries).Error
        if err == nil {
            if err = json.Unmarshal(entries.Data, &contentData); err != nil {
                return contract.ArticleDetail{}, err
            }
        } else if !errors.Is(err, gorm.ErrRecordNotFound) {
            return contract.ArticleDetail{}, err
        }
    }
    
    return contract.ArticleDetail{
        ID:                 article.ID,
        CreatedAt:          article.CreatedAt,
        UpdatedAt:          article.UpdatedAt,
        ContentTypeID:      cloneUint64(article.ContentTypeID),
        AuthorID:           article.AuthorID,
        URL:                cloneString(article.URL),
        Slug:               cloneString(article.Slug),
        Status:             domain.Status(article.Status),
        CurrentVersionID:   *article.CurrentVersionID,
        PublishedVersionID: cloneUint64(article.PublishedVersionID),
        CurrentRevision: contract.RevisionDetail{
            ID:        revision.ID,
            VersionNo: revision.VersionNo,
            Title:     revision.Title,
            Content:   revision.Content,
            Summary:   revision.Summary,
            ChangeLog: revision.ChangeLog,
            CreatedAt: revision.CreatedAt,
        },
        CategoryIDs:    categoryIDs,
        Tags:           tags,
        PublishedAt:    publishedAt,
        ContentSchemas: contentSchemas,
        ContentData:    contentData,
    }, nil
}

type articleRecord struct {
    ID                 uint64
    CreatedAt          time.Time
    UpdatedAt          time.Time
    DeletedAt          gorm.DeletedAt
    CurrentVersionID   *uint64
    PublishedVersionID *uint64
    ContentTypeID      *uint64
    AuthorID           uint64
    URL                *string
    Slug               *string
    Title              string
    Summary            string
    Status             int8
}

func (articleRecord) TableName() string {
    return "article"
}

type articleVersionRecord struct {
    ID        uint64
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt
    ArticleID uint64
    VersionNo uint64
    Title     string
    Content   string
    Summary   string
    ChangeLog string
}

func (articleVersionRecord) TableName() string {
    return "article_version"
}

type articlePublishRecord struct {
    ID          uint64
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   gorm.DeletedAt
    ArticleID   uint64
    VersionID   uint64
    PublishAt   time.Time
    PublishType int8
    Remark      string
}

func (articlePublishRecord) TableName() string {
    return "article_publish"
}

type articleCategoryRecord struct {
    ArticleID  uint64
    CategoryID uint64
}

func (articleCategoryRecord) TableName() string {
    return "article_category_rel"
}

type articleTagRecord struct {
    ID        uint64
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt
    Name      string
}

func (articleTagRecord) TableName() string {
    return "article_tag"
}

type articleTagRelationRecord struct {
    ArticleID uint64
    TagID     uint64
}

func (articleTagRelationRecord) TableName() string {
    return "article_tag_rel"
}

func cloneUint64(value *uint64) *uint64 {
    if value == nil {
        return nil
    }
    result := *value
    return &result
}

func cloneString(value *string) *string {
    if value == nil {
        return nil
    }
    result := *value
    return &result
}

func boolValue(value *bool, fallback bool) bool {
    if value == nil {
        return fallback
    }
    return *value
}
