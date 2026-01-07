package errno

import `dpcms/internal/erroz`

var (
    UserNotExists          = erroz.New(erroz.Code(erroz.ModuleUser, erroz.TypNotFound, 0), "用户不存在")
    UserNameNotExists      = erroz.New(erroz.Code(erroz.ModuleUser, erroz.TypNotFound, 1), "用户名不存在")
    UserNameExists         = erroz.New(erroz.Code(erroz.ModuleUser, erroz.TypConflict, 0), "用户名已存在")
    UserEmailExists        = erroz.New(erroz.Code(erroz.ModuleUser, erroz.TypConflict, 1), "邮箱已存在")
    UserWrongPasswd        = erroz.New(erroz.Code(erroz.ModuleUser, erroz.TypParameter, 0), "密码错误")
    UserWrongConfirmPasswd = erroz.New(erroz.Code(erroz.ModuleUser, erroz.TypParameter, 1), "两次密码输入不一致")
    UserEqualsOldPasswd    = erroz.New(erroz.Code(erroz.ModuleUser, erroz.TypParameter, 2), "新密码不能和旧密码相同")
)
