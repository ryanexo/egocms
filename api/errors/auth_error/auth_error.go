package auth_error

import `dpcms/erroz`

var (
    ErrUnauthorized         = erroz.New("auth.unauthorized", "未授权")
    ErrAuthorizationExpired = erroz.New("auth.expired", "授权已过期，请重新登陆")
)
