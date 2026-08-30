package httpx

import `cms/internal/public/erroz`

var (
    OK               = erroz.NewError("OK", "操作成功")
    Unknown          = erroz.NewError("SERV_001", "系统异常")
    DataNotFound     = erroz.NewError("SERV_002", "数据不存在")
    ValidationFailed = erroz.NewError("SERV_003", "参数错误")
)
