package dto

type ApiResult[T any] struct {
    Data T      `json:"data,omitempty"`
    Code string `json:"code"`
    Msg  string `json:"msg"`
}

type ApiPageData[T any] struct {
    Total    int64 `json:"total"`
    PageSize int   `json:"pageSize"`
    PageNo   int   `json:"pageNo"`
    List     []T   `json:"list"`
}

type ApiPagedResult[T any] struct {
    Data ApiPageData[T] `json:"data"`
    Code string         `json:"code"`
    Msg  string         `json:"msg"`
}

type ApiArticle = ApiResult[Article]
type ApiArticleList = ApiPagedResult[Article]
