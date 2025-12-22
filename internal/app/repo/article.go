package repo

import (
    `context`
    
    `dpcms/internal/infra/persistence/model`
    `dpcms/internal/infra/persistence/query`
    
    `gorm.io/gen/field`
)

type Article struct {
    persist *query.Query
}

func NewArticle(persist *query.Query) *Article {
    return &Article{persist: persist}
}

func (r Article) Find(ctx context.Context, id uint64) (*model.Article, error) {
    return r.persist.Article.WithContext(ctx).Preload(field.Associations).Where(r.persist.Article.ID.Eq(id)).First()
}

func (r Article) Create(ctx context.Context, article *model.Article) error {
    return r.persist.Article.WithContext(ctx).Create(article)
}
