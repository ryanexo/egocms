package menu_error

import (
    `dpcms/internal/erroz`
)

var (
    ErrCircularReferenceWhenMove = erroz.New("menu.circular", "目标菜单 %s 为当前菜单 %s 的子级")
)
