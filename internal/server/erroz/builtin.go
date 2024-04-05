package erroz

var (
    OK     = New("ok", "操作成功")
    Failed = New("server.error", "服务器错误")

    ErrValidation       = New("client.param_error", "参数错误")
    ErrAccountNotExists = New("account.not_exists", "帐号不存在")
    ErrUnameNotExists   = New("account.username_not_exists", "用户名不存在")
    ErrMailNotExists    = New("account.email_not_exists", "邮箱不存在")
    ErrUnameExists      = New("account.username_exists", "用户名已存在")
    ErrMailExists       = New("account.mail_exists", "邮箱已存在")
    ErrWrongPassword    = New("account.passwd_wrong", "密码错误")
    ErrRePasswordNotEq  = New("account.repassword_not_eq", "两次密码输入不一致")

    ErrUnauthorized = New("auth.unauthorized", "未授权")
)
