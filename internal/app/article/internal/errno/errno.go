package errno

import (
    `cms/internal/erroz`
    `cms/internal/httpserver`
)

const (
    statusCannotTransform = erroz.Article + iota
)

var (
    ErrInvalidStatusTransition = httpserver.NewError(statusCannotTransform, "当前状态不允许此操作")
)
