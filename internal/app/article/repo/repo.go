package repo

import (
    `context`
    
    `dpcms/internal/infra/persistence/datatype`
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/infra/persistence/query`
    
    `gorm.io/gen`
)

type ArticleRepo struct {
    query *query.Query
}

func NewArticleRepo(persist *query.Query) *ArticleRepo {
    return &ArticleRepo{persist}
}

func (r *ArticleRepo) Create(ctx context.Context, article *model.Article) error {
    return r.query.Article.WithContext(ctx).Create(article)
}

func (r *ArticleRepo) FindByID(ctx context.Context, id datatype.SafeUint64) (*model.Article, error) {
    artModel := r.query.Article
    schemaModel := r.query.ArticleModelSchema
    return artModel.WithContext(ctx).
        Preload(
            artModel.Keywords.Select(r.query.ArticleKeywords.Keyword),
            artModel.ModelData.Select(r.query.ArticleModelJsonData.Data),
            artModel.ModelSchema.Select(
                schemaModel.Type,
                schemaModel.FieldName,
                schemaModel.FieldKey,
                schemaModel.Description,
            ),
            artModel.Author.Select(r.query.UserProfile.Nickname),
            artModel.Category.Select(r.query.Category.Name),
        ).
        Where(artModel.ID.Eq(id.Raw())).
        First()
}

func (r *ArticleRepo) FindByIDWithContent(ctx context.Context, id datatype.SafeUint64) (*model.Article, error) {
    artData, err := r.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    content, err := r.query.ArticleContent.WithContext(ctx).Where(r.query.ArticleContent.ArticleID.Eq(id.Raw())).First()
    if err != nil {
        return nil, err
    }
    artData.Content = content
    return artData, nil
}

func (r *ArticleRepo) FindByIDWithoutPreload(ctx context.Context, id datatype.SafeUint64) (*model.Article, error) {
    dao := r.query.Article
    return dao.WithContext(ctx).Where(dao.ID.Eq(id.Raw())).First()
}

func (r *ArticleRepo) Update(ctx context.Context, article *model.Article) error {
    _, err := r.query.Article.WithContext(ctx).Where(r.query.Article.ID.Eq(article.ID.Raw())).Updates(article)
    return err
}

func (r *ArticleRepo) UpdateStatus(ctx context.Context, id datatype.SafeUint64, status int8) (gen.ResultInfo, error) {
    dao := r.query.Article
    return dao.WithContext(ctx).Where(dao.ID.Eq(id.Raw())).Update(dao.Status, status)
}

func (r *ArticleRepo) UpdateContent(ctx context.Context, id datatype.SafeUint64, content string) error {
    contentModel := r.query.ArticleContent
    _, err := contentModel.WithContext(ctx).Where(contentModel.ArticleID.Eq(id.Raw())).Update(contentModel.Content, content)
    return err
}

func (r *ArticleRepo) ReplaceKeywords(ctx context.Context, id datatype.SafeUint64, keywords []string) error {
    err := r.DeleteKeywords(ctx, id)
    if err != nil {
        return err
    }
    
    if len(keywords) > 0 {
        keywordSlice := make([]*model.ArticleKeywords, 0, len(keywords))
        for i := range keywords {
            keywordSlice[i] = &model.ArticleKeywords{ArticleID: id, Keyword: keywords[i]}
        }
        return r.query.ArticleKeywords.WithContext(ctx).Create(keywordSlice...)
    }
    
    return nil
}

func (r *ArticleRepo) DeleteArticle(ctx context.Context, id datatype.SafeUint64) error {
    _, err := r.query.Article.WithContext(ctx).Where(r.query.Article.ID.Eq(id.Raw())).Delete()
    return err
}

func (r *ArticleRepo) DeleteContent(ctx context.Context, id datatype.SafeUint64) error {
    _, err := r.query.ArticleContent.WithContext(ctx).Where(r.query.ArticleContent.ArticleID.Eq(id.Raw())).Delete()
    return err
}

func (r *ArticleRepo) DeleteKeywords(ctx context.Context, id datatype.SafeUint64) error {
    _, err := r.query.ArticleKeywords.WithContext(ctx).Unscoped().Where(r.query.ArticleKeywords.ArticleID.Eq(id.Raw())).Delete()
    return err
}
