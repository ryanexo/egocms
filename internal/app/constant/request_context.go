package constant

type requestContextKey string

func (s requestContextKey) Create(str string) string {
    return string(s) + "." + str
}

const uni = requestContextKey("uni")

var (
    RequestUserKey = uni.Create("CurrentUser")
)
