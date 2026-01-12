package errno

import (
    `dpcms/internal/erroz`
    `dpcms/internal/erroz/type`
)

var (
    FileDriverConfigNotExists  = erroz.New(erroz.Code("FILE", errtype.NotFound, 0), "未配置文件驱动")
    FileDriverNotExists        = erroz.New(erroz.Code("FILE", errtype.NotFound, 1), "文件驱动不存在")
    FileDeleteFailedInCreating = erroz.New(erroz.Code("FILE", errtype.Unknown, 2), "上传失败")
    FilePreCreateFileFailed    = erroz.New(erroz.Code("FILE", errtype.Unknown, 3), "上传失败")
)
