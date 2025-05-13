package erroz

var (
    OK = New("ok", "操作成功")
    
    ErrUnknown = New("server.error", "服务器错误")
    
    ErrValidation = New("client.param.error", "参数错误")
    
    ErrUserIDNotExists   = New("user.not_exists.id", "帐号不存在")
    ErrUsernameNotExists = New("user.not_exists.username", "用户名不存在")
    ErrUsernameExists    = New("user.exists.username", "用户名已存在")
    ErrEMailExists       = New("user.exists.email", "邮箱已存在")
    ErrWrongPassword     = New("user.password.wrong", "密码错误")
    
    ErrUnauthorized = New("auth.unauthorized", "未授权")
    
    ErrAuthorizationExpired = New("auth.expired", "授权已过期，请重新登陆")
    
    ErrTargetNodeIsSourceChild = New("tree.target_is_source_child", "目标节点为源节点的子节点")
)
