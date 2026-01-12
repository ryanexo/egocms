package errno

import (
    `dpcms/internal/erroz`
    `dpcms/internal/erroz/type`
)

var (
    RoleCircular = erroz.New(erroz.Code("ROLE", errtype.Conflict, 0), "角色 %s 已继承自当前角色，无法建立继承关系")
)
