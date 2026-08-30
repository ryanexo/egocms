package authz

import "cms/internal/public/erroz"

var (
    ErrAuthorized   = erroz.NewError("AUTH_001", "未授权")
    ErrAccessDenied = erroz.NewError("AUTH_002", "无权访问")
)
