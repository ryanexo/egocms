package errno

import `dpcms/internal/erroz`

var (
    Unauthorized         = erroz.New(erroz.Code(erroz.ModuleToken, erroz.TypAuth, 0), "未授权")
    AuthorizationExpired = erroz.New(erroz.Code(erroz.ModuleServer, erroz.TypAuth, 1), "授权已过期，请重新登陆")
)
