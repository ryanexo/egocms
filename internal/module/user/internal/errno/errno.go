package errno

import (
    `cms/internal/erroz`
    `cms/internal/erroz/type`
)

var (
    UserNotExists          = erroz.New(erroz.Code("USER", errtype.NotFound, 0), "用户不存在")
    UserNameNotExists      = erroz.New(erroz.Code("USER", errtype.NotFound, 1), "用户名不存在")
    UserNameExists         = erroz.New(erroz.Code("USER", errtype.Conflict, 0), "用户名已存在")
    UserEmailExists        = erroz.New(erroz.Code("USER", errtype.Conflict, 1), "邮箱已存在")
    UserWrongPasswd        = erroz.New(erroz.Code("USER", errtype.Parameter, 0), "密码错误")
    UserWrongConfirmPasswd = erroz.New(erroz.Code("USER", errtype.Parameter, 1), "两次密码输入不一致")
    UserEqualsOldPasswd    = erroz.New(erroz.Code("USER", errtype.Parameter, 2), "新密码不能和旧密码相同")
)
