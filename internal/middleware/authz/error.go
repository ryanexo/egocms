package authz

import (
    `dpcms/internal/erroz`
    errtype `dpcms/internal/erroz/type`
)

var ErrAuthorized = erroz.New(erroz.Code("AUTH", errtype.Auth, 0), "未授权")
var ErrAccessDenied = erroz.New(erroz.Code("AUTH", errtype.Auth, 1), "无权访问")
