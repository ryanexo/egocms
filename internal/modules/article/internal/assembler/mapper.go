package assembler

import (
    `database/sql`
    `time`
    
    `cms/internal/infra/store/model`
    dto2 `cms/internal/modules/article/internal/dto`
    `cms/internal/public/apitype`
    model3 `cms/internal/public/model`
)

func extractKeywords(data []*model3.ArticleTag) []string {
    keywords := make([]string, 0, len(data))
    if len(data) > 0 {
        for _, item := range data {
            keywords = append(keywords, item.Keyword)
        }
    }
    return keywords
}

func ToArticleCreateCommand(user *model3.User, data *dto2.ArticleCreateParams) *model3.Article {
    result := &model3.Article{
        AuthorID:       user.ID,
        Url:            data.Url,
        CategoryID:     data.CategoryID,
        ContentTypeID:  data.ModelID,
        Flag:           data.Flag,
        Title:          data.Title,
        Summary:        data.Description,
        Target:         sql.NullString{String: data.Target, Valid: true},
        Content:        &model.ArticleContent{Content: data.Content},
        ContentEntries: &model3.ContentTypeEntries{},
    }
    
    keywordCount := len(data.Keywords)
    if keywordCount > 0 {
        keywords := make([]*model3.ArticleTag, 0, len(data.Keywords))
        for _, keyword := range data.Keywords {
            keywords = append(keywords, &model3.ArticleTag{Keyword: keyword})
        }
        result.Keywords = keywords
    }
    
    return result
}

func ToArticleDTO(data *model3.Article) *dto2.Article {
    result := &dto2.Article{
        Base: apitype.Base{
            ID:        data.ID,
            CreatedAt: data.CreatedAt,
            UpdatedAt: data.UpdatedAt,
        },
        Url:         data.Url,
        CategoryID:  data.CategoryID,
        AuthorID:    data.AuthorID,
        Flag:        data.Flag,
        Title:       data.Title,
        Description: data.Summary,
        ClickCount:  data.ClickCount,
        Status:      data.Status,
        Target:      data.Target.String,
        Keywords:    extractKeywords(data.Keywords),
        ModelID:     data.ContentTypeID,
        PublishAt:   time.Time{},
    }
    
    if data.Content != nil {
        result.Content = &data.Content.Content
    }
    
    if data.Category != nil {
        result.CategoryName = data.Category.Name
    }
    
    if data.Author != nil {
        result.AuthorName = data.Author.Nickname
    }
    
    if data.ContentTypeSchemas != nil {
        schema := make([]*dto2.ArticleModelSchema, 0, len(data.ContentTypeSchemas))
        for _, modelSchema := range data.ContentTypeSchemas {
            schema = append(schema, &dto2.ArticleModelSchema{
                FieldKey:  modelSchema.FieldKey,
                FieldName: modelSchema.FieldName,
                Type:      modelSchema.Type,
                Sequence:  modelSchema.Sequence,
            })
        }
    }
    
    if data.ContentEntries != nil || data.ContentEntries.Data != nil {
        result.ModelData = data.ContentEntries.Data
    }
    
    return result
}
