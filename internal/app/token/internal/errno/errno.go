package errno

import (
    `cms/internal/erroz`
    erroz2 `cms/internal/httpserver`
)

var (
    Unauthorized         = erroz2.New(erroz.Code("TOKEN", errtype.Auth, 0), "未授权")
    AuthorizationExpired = erroz2.New(erroz.Code("TOKEN", errtype.Auth, 1), "授权已过期，请重新登陆")
)
