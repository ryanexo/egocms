package authz

import (
    `cms/internal/httpserver`
)

var (
    ErrAuthorized   = httpserver.NewError(400001, "未授权")
    ErrAccessDenied = httpserver.NewError(400002, "无权访问")
)
