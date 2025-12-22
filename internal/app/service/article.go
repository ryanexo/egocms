package service

import (
    `context`
    
    `dpcms/internal/app/dto`
    `dpcms/internal/app/repo`
)

type Article struct {
    repo *repo.Repo
}

func NewArticle(repo *repo.Repo) *Article {
    return &Article{repo}
}

func (s Article) Create(ctx context.Context, params dto.ArticleCreateParams) {
}
