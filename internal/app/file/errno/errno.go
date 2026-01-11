package errno

import `dpcms/internal/erroz`

var (
    FileDriverConfigNotExists  = erroz.New(erroz.Code(erroz.ModuleFile, erroz.TypNotFound, 0), "未配置文件驱动")
    FileDriverNotExists        = erroz.New(erroz.Code(erroz.ModuleFile, erroz.TypNotFound, 1), "文件驱动不存在")
    FileDeleteFailedInCreating = erroz.New(erroz.Code(erroz.ModuleFile, erroz.TypUnknown, 2), "上传失败，已上传文件回滚失败")
    FileGeneratePathFailed     = erroz.New(erroz.Code(erroz.ModuleFile, erroz.TypUnknown, 3), "上传失败")
)
