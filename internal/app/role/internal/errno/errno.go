package errno

import (
    `cms/internal/erroz`
    erroz2 `cms/internal/httpserver`
)

var (
    RoleCircular = erroz2.New(erroz.Code("ROLE", errtype.Conflict, 0), "角色 %s 已继承自当前角色，无法建立继承关系")
)
