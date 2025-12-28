package article

import (
    `dpcms/internal/app/dto`
    `dpcms/internal/infra/persistence/customvalue`
    `dpcms/internal/infra/persistence/model`
)

type DataViewContext struct {
    Article   *model.Article
    Category  *model.Category
    Author    *model.UserProfile
    Keywords  []*model.ArticleKeywords
    Schema    []*model.ArticleModelSchema
    ModelData map[string]any
}

func NewArticleBasicData(ctx DataViewContext) *dto.Article {
    keywords := make([]string, len(ctx.Keywords))
    schema := make([]dto.ArticleModelSchema, len(ctx.Schema))
    
    for _, keyword := range ctx.Keywords {
        keywords = append(keywords, keyword.Keyword)
    }
    
    for _, modelSchema := range ctx.Schema {
        schema = append(schema, dto.ArticleModelSchema{
            FieldKey:  modelSchema.FieldKey,
            FieldName: modelSchema.FieldName,
            Type:      modelSchema.Type,
            Sequence:  modelSchema.Sequence,
        })
    }
    
    return &dto.Article{
        Base:         ctx.Article.Base,
        Url:          ctx.Article.Url,
        CategoryId:   customvalue.NewUint64String(ctx.Category.ID),
        CategoryName: ctx.Category.Name,
        AuthorId:     customvalue.NewUint64String(ctx.Author.ID),
        AuthorName:   ctx.Author.Nickname,
        Flag:         ctx.Article.Flag,
        Title:        ctx.Article.Title,
        Description:  ctx.Article.Description,
        ClickCount:   customvalue.NewUint64String(ctx.Article.ClickCount),
        Status:       ctx.Article.Status,
        Target:       ctx.Article.Target.String,
        Keywords:     keywords,
        ModelId:      customvalue.NewUint64String(ctx.Article.ModelId),
        ModelSchema:  schema,
        ModelData:    ctx.ModelData,
    }
}
