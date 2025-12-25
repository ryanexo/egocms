package repo

import (
    `context`
    
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/infra/persistence/query`
    
    `gorm.io/gen/field`
)

type Article struct {
    context context.Context
    persist *query.Query
}

func NewArticle(ctx context.Context, persist *query.Query) *Article {
    return &Article{ctx, persist}
}

func (r *Article) Find(id uint64) (*model.Article, error) {
    return r.persist.Article.WithContext(r.context).Preload(field.Associations).Where(r.persist.Article.ID.Eq(id)).First()
}

func (r *Article) Create(article *model.Article) error {
    return r.persist.Article.WithContext(r.context).Create(article)
}

func (r *Article) UpdateKeywords(id uint64, keywords []string) error {
    _, err := r.persist.ArticleKeywords.WithContext(r.context).Where(r.persist.ArticleKeywords.ArticleID.Eq(id)).Delete()
    if err != nil {
        return err
    }
    keywordSlice := make([]*model.ArticleKeywords, 0, len(keywords))
    for i := range keywords {
        keywordSlice[i] = &model.ArticleKeywords{ArticleID: id, Keyword: keywords[i]}
    }
    return r.persist.ArticleKeywords.WithContext(r.context).Create(keywordSlice...)
}

func (r *Article) SaveContentModelData(id uint64, data []*model.ArticleModelData) error {
    defDao := r.persist.ContentModelData
    _, err := defDao.WithContext(r.context).Where(defDao.ContentModelId.Eq(id)).Delete()
    if err != nil {
        return err
    }
    return defDao.WithContext(r.context).Create(data...)
}
