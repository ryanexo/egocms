package domain

import (
    `cms/internal/public/erroz`
)

var (
    ErrTitleRequired    = erroz.NewError(10001, "标题必填")
    ErrTitleTooLong     = erroz.NewError(10002, "标题长度超出限制")
    ErrSummaryTooLong   = erroz.NewError(10003, "简介长度超出限制")
    ErrContentTooLong   = erroz.NewError(10004, "内容长度超出限制")
    ErrChangeLogTooLong = erroz.NewError(10005, "修改日志长度超出限制")
)
