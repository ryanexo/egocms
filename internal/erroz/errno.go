package erroz

import `cms/internal/erroz/type`

var (
    OK               = New(Code("SERVER", errtype.OK, 0), "操作成功")
    Unknown          = New(Code("SERVER", errtype.Unknown, 0), "系统异常")
    DataNotFound     = New(Code("SERVER", errtype.NotFound, 0), "数据不存在")
    ValidationFailed = New(Code("CLIENT", errtype.Parameter, 0), "参数错误")
)
