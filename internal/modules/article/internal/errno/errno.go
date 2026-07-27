package errno

import (
    `cms/internal/erroz`
    `cms/internal/httpx`
)

const (
    statusCannotTransform = erroz.Article + iota
)

var (
    ErrInvalidStatusTransition = httpx.NewError(statusCannotTransform, "当前状态不允许此操作")
)
