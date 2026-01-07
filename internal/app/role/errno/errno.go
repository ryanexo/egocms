package errno

import `dpcms/internal/erroz`

var (
    RoleCircular = erroz.New(erroz.Code(erroz.ModuleRole, erroz.TypConflict, 0), "角色 %s 已继承自当前角色，无法建立继承关系")
)
