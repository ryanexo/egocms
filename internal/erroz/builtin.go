package erroz

var (
    OK              = New("ok", "操作成功")
    ErrUnknown      = New("server.error", "服务器错误")
    ErrDataNotFound = New("server.data.notfound", "数据不存在")
    ErrValidation   = New("client.param.error", "参数错误")
)
