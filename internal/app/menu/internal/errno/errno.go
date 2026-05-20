package errno

import (
    `cms/internal/erroz`
    erroz2 `cms/internal/httpserver`
)

var (
    MenuCircular = erroz2.New(erroz.Code("MENU", errtype.Conflict, 0), "目标菜单 %s 为当前菜单 %s 的子级，无法移动")
)
