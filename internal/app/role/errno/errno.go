package errno

import "cms/internal/public/erroz"

var (
    ErrNameRequired             = erroz.NewError("ROLE_001", "角色名称必填")
    ErrInvalidNameLength        = erroz.NewError("ROLE_002", "角色名称长度超出限制")
    ErrDescriptionRequired      = erroz.NewError("ROLE_003", "角色描述必填")
    ErrInvalidDescriptionLength = erroz.NewError("ROLE_004", "角色描述长度超出限制")
    ErrCircular                 = erroz.NewError("ROLE_005", "角色 %s 已继承当前角色，无法建立继承关系")
    ErrTooManyInheritedRoles    = erroz.NewError("ROLE_006", "角色最多继承 10 个角色")
    ErrInheritSelf              = erroz.NewError("ROLE_007", "不能继承自己")
)
