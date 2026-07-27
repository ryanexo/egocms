package authz

import (
    `cms/internal/httpx`
)

var (
    ErrAuthorized   = httpx.NewError(400001, "未授权")
    ErrAccessDenied = httpx.NewError(400002, "无权访问")
)
