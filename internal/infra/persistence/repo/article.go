package repo

import (
    `context`
    
    `dpcms/internal/infra/persistence/datatype`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/infra/persistence/query`
)

type Article struct {
    persist *query.Query
}

func NewArticle(persist *query.Query) Article {
    return Article{persist}
}

func (r Article) Create(ctx context.Context, article *model.Article) error {
    return r.persist.Article.WithContext(ctx).Create(article)
}

func (r Article) FindByID(ctx context.Context, id datatype.SafeUint64) (*model.Article, error) {
    artModel := r.persist.Article
    schemaModel := r.persist.ArticleModelSchema
    return artModel.WithContext(ctx).
        Preload(
            artModel.Keywords.Select(r.persist.ArticleKeywords.Keyword),
            artModel.ModelData.Select(r.persist.ArticleModelJsonData.Data),
            artModel.ModelSchema.Select(
                schemaModel.Type,
                schemaModel.FieldName,
                schemaModel.FieldKey,
                schemaModel.Description,
            ),
            artModel.Author.Select(r.persist.UserProfile.Nickname),
            artModel.Category.Select(r.persist.Category.Name),
        ).
        Where(artModel.ID.Eq(id.Raw())).
        First()
}

func (r Article) FindByIDWithContent(ctx context.Context, id datatype.SafeUint64) (*model.Article, error) {
    artData, err := r.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    content, err := r.persist.ArticleContent.WithContext(ctx).Where(r.persist.ArticleContent.ArticleID.Eq(id.Raw())).First()
    if err != nil {
        return nil, err
    }
    artData.Content = content
    return artData, nil
}

func (r Article) Update(ctx context.Context, article *model.Article) error {
    _, err := r.persist.Article.WithContext(ctx).Where(r.persist.Article.ID.Eq(article.ID.Raw())).Updates(article)
    return err
}

func (r Article) UpdateContent(ctx context.Context, id datatype.SafeUint64, content string) error {
    contentModel := r.persist.ArticleContent
    _, err := contentModel.WithContext(ctx).Where(contentModel.ArticleID.Eq(id.Raw())).Update(contentModel.Content, content)
    return err
}

func (r Article) ReplaceKeywords(ctx context.Context, id datatype.SafeUint64, keywords []string) error {
    err := r.DeleteKeywords(ctx, id)
    if err != nil {
        return err
    }
    
    if len(keywords) > 0 {
        keywordSlice := make([]*model.ArticleKeywords, 0, len(keywords))
        for i := range keywords {
            keywordSlice[i] = &model.ArticleKeywords{ArticleID: id, Keyword: keywords[i]}
        }
        return r.persist.ArticleKeywords.WithContext(ctx).Create(keywordSlice...)
    }
    
    return nil
}

func (r Article) DeleteArticle(ctx context.Context, id datatype.SafeUint64) error {
    _, err := r.persist.Article.WithContext(ctx).Where(r.persist.Article.ID.Eq(id.Raw())).Delete()
    return err
}

func (r Article) DeleteContent(ctx context.Context, id datatype.SafeUint64) error {
    _, err := r.persist.ArticleContent.WithContext(ctx).Where(r.persist.ArticleContent.ArticleID.Eq(id.Raw())).Delete()
    return err
}

func (r Article) DeleteKeywords(ctx context.Context, id datatype.SafeUint64) error {
    _, err := r.persist.ArticleKeywords.WithContext(ctx).Unscoped().Where(r.persist.ArticleKeywords.ArticleID.Eq(id.Raw())).Delete()
    return err
}
