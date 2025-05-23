package role_error

import (
    `dpcms/erroz`
)

var (
    ErrRoleCircularReference = erroz.New("role.circular", "角色 %s 已继承自当前角色 %s，无法建立继承关系")
)
