package adapter

import (
    `context`
    
    `cms/internal/app/article/contract`
    `cms/internal/infra/persist/datatype`
    `cms/internal/infra/persist/model`
    `cms/internal/infra/persist/query`
    
    `gorm.io/gen`
)

type articleRepo struct {
    query *query.Query
}

var _ contract.ArticleRepo = (*articleRepo)(nil)

func NewArticleRepo(persist *query.Query) contract.ArticleRepo {
    return articleRepo{persist}
}

func (s articleRepo) CloneWithQuery(q *query.Query) contract.ArticleRepo {
    return NewArticleRepo(q)
}

func (s articleRepo) Create(ctx context.Context, article *model.Article) error {
    return s.query.Article.WithContext(ctx).Create(article)
}

func (s articleRepo) FindByID(ctx context.Context, id datatype.SafeUint64) (*model.Article, error) {
    artModel := s.query.Article
    schemaModel := s.query.ArticleModelSchema
    return artModel.WithContext(ctx).
        Preload(
            artModel.Keywords.Select(s.query.ArticleKeywords.Keyword),
            artModel.ModelData.Select(s.query.ArticleModelJsonData.Data),
            artModel.ModelSchema.Select(
                schemaModel.Type,
                schemaModel.FieldName,
                schemaModel.FieldKey,
                schemaModel.Description,
            ),
            artModel.Author.Select(s.query.UserProfile.Nickname),
            artModel.Category.Select(s.query.Category.Name),
        ).
        Where(artModel.ID.Eq(id.Raw())).
        First()
}

func (s articleRepo) FindByIDWithContent(ctx context.Context, id datatype.SafeUint64) (*model.Article, error) {
    artData, err := s.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    content, err := s.query.ArticleContent.WithContext(ctx).Where(s.query.ArticleContent.ArticleID.Eq(id.Raw())).First()
    if err != nil {
        return nil, err
    }
    artData.Content = content
    return artData, nil
}

func (s articleRepo) FindByIDWithoutPreload(ctx context.Context, id datatype.SafeUint64) (*model.Article, error) {
    dao := s.query.Article
    return dao.WithContext(ctx).Where(dao.ID.Eq(id.Raw())).First()
}

func (s articleRepo) Update(ctx context.Context, article *model.Article) error {
    _, err := s.query.Article.WithContext(ctx).Where(s.query.Article.ID.Eq(article.ID.Raw())).Updates(article)
    return err
}

func (s articleRepo) UpdateStatus(ctx context.Context, id datatype.SafeUint64, status int8) (gen.ResultInfo, error) {
    dao := s.query.Article
    return dao.WithContext(ctx).Where(dao.ID.Eq(id.Raw())).Update(dao.Status, status)
}

func (s articleRepo) UpdateContent(ctx context.Context, id datatype.SafeUint64, content string) error {
    contentModel := s.query.ArticleContent
    _, err := contentModel.WithContext(ctx).Where(contentModel.ArticleID.Eq(id.Raw())).Update(contentModel.Content, content)
    return err
}

func (s articleRepo) ReplaceKeywords(ctx context.Context, id datatype.SafeUint64, keywords []string) error {
    err := s.DeleteKeywords(ctx, id)
    if err != nil {
        return err
    }
    
    if len(keywords) > 0 {
        keywordSlice := make([]*model.ArticleKeywords, 0, len(keywords))
        for i := range keywords {
            keywordSlice[i] = &model.ArticleKeywords{ArticleID: id, Keyword: keywords[i]}
        }
        return s.query.ArticleKeywords.WithContext(ctx).Create(keywordSlice...)
    }
    
    return nil
}

func (s articleRepo) DeleteArticle(ctx context.Context, id datatype.SafeUint64) error {
    _, err := s.query.Article.WithContext(ctx).Where(s.query.Article.ID.Eq(id.Raw())).Delete()
    return err
}

func (s articleRepo) DeleteContent(ctx context.Context, id datatype.SafeUint64) error {
    _, err := s.query.ArticleContent.WithContext(ctx).Where(s.query.ArticleContent.ArticleID.Eq(id.Raw())).Delete()
    return err
}

func (s articleRepo) DeleteKeywords(ctx context.Context, id datatype.SafeUint64) error {
    _, err := s.query.ArticleKeywords.WithContext(ctx).Unscoped().Where(s.query.ArticleKeywords.ArticleID.Eq(id.Raw())).Delete()
    return err
}
