package errno

import "cms/internal/public/erroz"

var (
	ErrUserNotFound         = erroz.NewError("USER_001", "用户不存在")
	ErrUsernameExists       = erroz.NewError("USER_002", "用户名已存在")
	ErrEmailExists          = erroz.NewError("USER_003", "邮箱已存在")
	ErrWrongPassword        = erroz.NewError("USER_004", "密码错误")
	ErrPasswordMismatch     = erroz.NewError("USER_005", "两次密码输入不一致")
	ErrPasswordUnchanged    = erroz.NewError("USER_006", "新密码不能和旧密码相同")
	ErrUsernameRequired     = erroz.NewError("USER_007", "用户名必填")
	ErrInvalidUsername      = erroz.NewError("USER_008", "用户名只能包含字母和数字，长度为 4 到 255 个字符")
	ErrEmailRequired        = erroz.NewError("USER_009", "邮箱必填")
	ErrInvalidEmail         = erroz.NewError("USER_010", "邮箱格式不正确或长度超过限制")
	ErrPasswordRequired     = erroz.NewError("USER_011", "密码必填")
	ErrInvalidProfileLength = erroz.NewError("USER_012", "用户资料字段长度超过限制")
	ErrInvalidGender        = erroz.NewError("USER_013", "性别取值无效")
)
