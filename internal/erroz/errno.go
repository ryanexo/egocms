package erroz

import (
    `cms/internal/httpserver`
)

var (
    OK               = httpserver.NewError(0, "操作成功")
    Unknown          = httpserver.NewError(50000, "系统异常")
    DataNotFound     = httpserver.NewError(50001, "数据不存在")
    ValidationFailed = httpserver.NewError(50002, "参数错误")
)
