package erroz

var (
    OK              = New("ok", "操作成功")
    ErrUnknown      = New("httpserver.error", "服务器错误")
    ErrValidation   = New("client.param.error", "参数错误")
    ErrDataNotFound = New("client.data.notfound", "数据不存在")
)
