package errno

import (
    `cms/internal/erroz`
    `cms/internal/erroz/type`
)

var (
    MenuCircular = erroz.New(erroz.Code("MENU", errtype.Conflict, 0), "目标菜单 %s 为当前菜单 %s 的子级，无法移动")
)
