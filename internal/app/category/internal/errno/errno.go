package errno

import (
    `cms/internal/erroz`
    erroz2 `cms/internal/httpserver`
)

var (
    CategoryCircular = erroz2.New(erroz.Code("CATEGORY", errtype.Parameter, 0), "目标分类 %s 为当前分类 %s 的子级，无法移动")
)
