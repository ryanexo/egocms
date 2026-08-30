package apitype

type ApiCreateResult struct {
    Data string `json:"data" example:"123456"` // 数据 ID
    Code string `json:"code"`
    Msg  string `json:"msg"`
}

type ApiEmptyResult struct {
    Code string `json:"code"`
    Msg  string `json:"msg"`
}
