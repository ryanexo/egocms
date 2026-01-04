package article

import (
    `dpcms/internal/app/dto`
    `dpcms/internal/app/dto/type`
    `dpcms/internal/infra/persistence/model`
)

type DataViewContext struct {
    Article  *model.Article
    Category *model.Category
    Author   *model.UserProfile
    Schema   []*model.ArticleModelSchema
}

func NewArticleDataView(ctx DataViewContext) *dto.Article {
    keywords := make([]string, len(ctx.Article.Keywords))
    schema := make([]*dto.ArticleModelSchema, len(ctx.Schema))
    
    for _, keyword := range ctx.Article.Keywords {
        keywords = append(keywords, keyword.Keyword)
    }
    
    for _, modelSchema := range ctx.Schema {
        schema = append(schema, &dto.ArticleModelSchema{
            FieldKey:  modelSchema.FieldKey,
            FieldName: modelSchema.FieldName,
            Type:      modelSchema.Type,
            Sequence:  modelSchema.Sequence,
        })
    }
    
    return &dto.Article{
        Base: dtotype.Base{
            ID:        ctx.Article.Base.ID,
            CreatedAt: ctx.Article.Base.CreatedAt,
            UpdatedAt: ctx.Article.Base.UpdatedAt,
        },
        Url:          ctx.Article.Url,
        CategoryID:   ctx.Category.ID,
        CategoryName: ctx.Category.Name,
        AuthorID:     ctx.Author.ID,
        AuthorName:   ctx.Author.Nickname,
        Flag:         ctx.Article.Flag,
        Title:        ctx.Article.Title,
        Description:  ctx.Article.Description,
        ClickCount:   ctx.Article.ClickCount,
        Status:       ctx.Article.Status,
        Target:       ctx.Article.Target.String,
        Keywords:     keywords,
        ModelID:      ctx.Article.ModelID,
        ModelSchema:  schema,
        ModelData:    ctx.Article.ModelData.Data,
    }
}
