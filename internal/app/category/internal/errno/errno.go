package errno

import (
    `cms/internal/erroz`
    `cms/internal/erroz/type`
)

var (
    CategoryCircular = erroz.New(erroz.Code("CATEGORY", errtype.Parameter, 0), "目标分类 %s 为当前分类 %s 的子级，无法移动")
)
