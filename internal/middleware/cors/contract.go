package cors

import `context`

type Options interface {
    GetAllowOrigin(context.Context) string
    GetAllowMethods(context.Context) string
    GetAllowHeaders(context.Context) string
    GetAllowCredentials(context.Context) string
    GetExposeHeaders(context.Context) string
}
