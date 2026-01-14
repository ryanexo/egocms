package service

import (
    `context`
    
    `cms/internal/infra/persistence/contract`
    `cms/internal/infra/persistence/datatype`
    `cms/internal/infra/persistence/model`
    
    `gorm.io/gen`
)

type ArticleRepo interface {
    contract.Repository[ArticleRepo]
    Create(ctx context.Context, article *model.Article) error
    FindByID(ctx context.Context, id datatype.SafeUint64) (*model.Article, error)
    FindByIDWithContent(ctx context.Context, id datatype.SafeUint64) (*model.Article, error)
    FindByIDWithoutPreload(ctx context.Context, id datatype.SafeUint64) (*model.Article, error)
    Update(ctx context.Context, article *model.Article) error
    UpdateStatus(ctx context.Context, id datatype.SafeUint64, status int8) (gen.ResultInfo, error)
    UpdateContent(ctx context.Context, id datatype.SafeUint64, content string) error
    ReplaceKeywords(ctx context.Context, id datatype.SafeUint64, keywords []string) error
    DeleteArticle(ctx context.Context, id datatype.SafeUint64) error
    DeleteContent(ctx context.Context, id datatype.SafeUint64) error
    DeleteKeywords(ctx context.Context, id datatype.SafeUint64) error
}
