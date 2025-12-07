package constant

type requestContextKey string

func (s requestContextKey) Create(str string) string {
    return string(s) + "." + str
}

const uni = requestContextKey("uni")

var (
    ApiAuthCurrentUser = uni.Create("CurrentUser")
)
