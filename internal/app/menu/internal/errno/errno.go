package errno

import `dpcms/internal/erroz`

var (
    MenuCircular = erroz.New(erroz.Code(erroz.ModuleMenu, erroz.TypConflict, 0), "目标菜单 %s 为当前菜单 %s 的子级，无法移动")
)
