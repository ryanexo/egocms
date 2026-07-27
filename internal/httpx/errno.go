package httpx

var (
    OK               = NewError(0, "操作成功")
    Unknown          = NewError(50000, "系统异常")
    DataNotFound     = NewError(50001, "数据不存在")
    ValidationFailed = NewError(50002, "参数错误")
)
