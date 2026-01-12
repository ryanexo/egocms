package errno

import (
    `dpcms/internal/erroz`
    `dpcms/internal/erroz/type`
)

var (
    Unauthorized         = erroz.New(erroz.Code("TOKEN", errtype.Auth, 0), "未授权")
    AuthorizationExpired = erroz.New(erroz.Code("TOKEN", errtype.Auth, 1), "授权已过期，请重新登陆")
)
