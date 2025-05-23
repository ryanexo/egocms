package user_error

import `dpcms/erroz`

var (
    ErrUserIDNotExists   = erroz.New("user.not_exists.id", "用户不存在")
    ErrUsernameNotExists = erroz.New("user.not_exists.username", "用户名不存在")
    ErrUsernameExists    = erroz.New("user.exists.username", "用户名已存在")
    ErrEMailExists       = erroz.New("user.exists.email", "邮箱已存在")
    ErrWrongPassword     = erroz.New("user.password.wrong", "密码错误")
)
