package domain

import `cms/internal/httpx`

var (
    ErrTitleRequired    = httpx.NewError(10001, "标题必填")
    ErrTitleTooLong     = httpx.NewError(10002, "标题长度超出限制")
    ErrSummaryTooLong   = httpx.NewError(10003, "简介长度超出限制")
    ErrContentTooLong   = httpx.NewError(10004, "内容长度超出限制")
    ErrChangeLogTooLong = httpx.NewError(10005, "修改日志长度超出限制")
)
