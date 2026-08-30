package errno

import (
    `cms/internal/erroz`
    erroz2 `cms/internal/public/erroz`
)

const (
    statusCannotTransform = erroz.Article + iota
)

var (
    ErrInvalidStatusTransition = erroz2.NewError(statusCannotTransform, "当前状态不允许此操作")
)
