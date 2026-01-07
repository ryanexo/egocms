package erroz

var (
    OK           = New(Code(ModuleServer, TypOK, 0), "操作成功")
    Unknown      = New(Code(ModuleServer, TypUnknown, 0), "系统异常")
    DataNotFound = New(Code(ModuleServer, TypNotFound, 0), "数据不存在")
    
    ValidationFailed = New(Code(ModuleClient, TypParameter, 0), "参数错误")
    
    Unauthorized = New(Code(ModuleServer, TypNotFound, 0), "未授权")
)
