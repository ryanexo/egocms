package assembler

import (
    `database/sql`
    `time`
    
    dto2 `cms/internal/domain/article/internal/dto`
    `cms/internal/infra/persistence/model`
    `cms/internal/util/types`
)

func extractKeywords(data []*model.ArticleKeywords) []string {
    keywords := make([]string, 0, len(data))
    if len(data) > 0 {
        for _, item := range data {
            keywords = append(keywords, item.Keyword)
        }
    }
    return keywords
}

func ToArticleCreateCommand(user *model.User, data *dto2.ArticleCreateParams) *model.Article {
    result := &model.Article{
        AuthorID:    user.ID,
        Url:         data.Url,
        CategoryID:  data.CategoryID,
        ModelID:     data.ModelID,
        Flag:        data.Flag,
        Title:       data.Title,
        Description: data.Description,
        Target:      sql.NullString{String: data.Target, Valid: true},
        Content:     &model.ArticleContent{Content: data.Content},
        ModelData:   &model.ArticleModelJsonData{},
    }
    
    keywordCount := len(data.Keywords)
    if keywordCount > 0 {
        keywords := make([]*model.ArticleKeywords, 0, len(data.Keywords))
        for _, keyword := range data.Keywords {
            keywords = append(keywords, &model.ArticleKeywords{Keyword: keyword})
        }
        result.Keywords = keywords
    }
    
    return result
}

func ToArticleDTO(data *model.Article) *dto2.Article {
    result := &dto2.Article{
        Base: types.Base{
            ID:        data.ID,
            CreatedAt: data.CreatedAt,
            UpdatedAt: data.UpdatedAt,
        },
        Url:         data.Url,
        CategoryID:  data.CategoryID,
        AuthorID:    data.AuthorID,
        Flag:        data.Flag,
        Title:       data.Title,
        Description: data.Description,
        ClickCount:  data.ClickCount,
        Status:      data.Status,
        Target:      data.Target.String,
        Keywords:    extractKeywords(data.Keywords),
        ModelID:     data.ModelID,
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
    
    if data.ModelSchema != nil {
        schema := make([]*dto2.ArticleModelSchema, 0, len(data.ModelSchema))
        for _, modelSchema := range data.ModelSchema {
            schema = append(schema, &dto2.ArticleModelSchema{
                FieldKey:  modelSchema.FieldKey,
                FieldName: modelSchema.FieldName,
                Type:      modelSchema.Type,
                Sequence:  modelSchema.Sequence,
            })
        }
    }
    
    if data.ModelData != nil || data.ModelData.Data != nil {
        result.ModelData = data.ModelData.Data
    }
    
    return result
}
