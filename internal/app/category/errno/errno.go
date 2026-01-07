package errno

import `dpcms/internal/erroz`

var (
    CategoryCircular = erroz.New(erroz.Code(erroz.ModuleCategory, erroz.TypParameter, 0), "目标分类 %s 为当前分类 %s 的子级，无法移动")
)
