package dto

type ApiResult[T any] struct {
    Data T      `json:"data"`
    Code string `json:"code"`
    Msg  string `json:"msg"`
}

type ApiPaginatedData[T any] struct {
    Total    int64 `json:"total"`
    PageSize int   `json:"pageSize"`
    PageNo   int   `json:"pageNo"`
    List     []T   `json:"list"`
}

type ApiPaginatedResult[T any] struct {
    Data ApiPaginatedData[T] `json:"data"`
    Code string              `json:"code"`
    Msg  string              `json:"msg"`
}

type ApiArticleModel = ApiResult[ArticleModel]
type ApiArticleModelList = ApiPaginatedResult[ArticleModel]
