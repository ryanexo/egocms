package article

import (
    `database/sql`
    
    `dpcms/internal/app/dto`
    `dpcms/internal/infra/persistence/model`
    
    `github.com/jinzhu/copier`
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

func BuildArticleCreateCommand(user *model.User, data *dto.ArticleCreateParams) *model.Article {
    result := &model.Article{
        AuthorID:    user.ID,
        Url:         data.Url,
        CategoryID:  data.CategoryId,
        ModelID:     data.ModelId,
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

func BuildArticleDTO(data *model.Article) (*dto.Article, error) {
    result := &dto.Article{Keywords: extractKeywords(data.Keywords)}
    if err := copier.Copy(&result, data); err != nil {
        return nil, err
    }
    
    if data.Category != nil {
        result.CategoryName = data.Category.Name
    }
    
    if data.Author != nil {
        result.AuthorName = data.Author.Nickname
    }
    
    if data.ModelSchema != nil {
        schema := make([]*dto.ArticleModelSchema, 0, len(data.ModelSchema))
        for _, modelSchema := range data.ModelSchema {
            schema = append(schema, &dto.ArticleModelSchema{
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
    
    return result, nil
}
