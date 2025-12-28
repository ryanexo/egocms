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

func (r *Article) FindBasicById(article *model.Article) error {
}

func (r *Article) FindDetailById(article *model.Article) error {
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

func (r *Article) SaveContentModelData(art model.Article, jsonMap map[string]any, data []*model.ArticleModelData) error {
    jsonData := model.ArticleModelJsonData{
        ArticleId: art.ID,
        ModelId:   art.ModelId,
        Data:      jsonMap,
    }
    err := r.persist.ArticleModelJsonData.WithContext(r.context).Create(&jsonData)
    if err != nil {
        return err
    }
    return r.persist.ArticleModelData.WithContext(r.context).Create(data...)
}
