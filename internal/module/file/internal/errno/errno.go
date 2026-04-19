package errno

import (
    `cms/internal/erroz`
    `cms/internal/erroz/type`
)

var (
    FileDriverConfigNotExists  = erroz.New(erroz.Code("FILE", errtype.NotFound, 0), "未配置文件驱动")
    FileDeleteFailedInCreating = erroz.New(erroz.Code("FILE", errtype.Unknown, 1), "上传失败")
    FilePreCreateFileFailed    = erroz.New(erroz.Code("FILE", errtype.Unknown, 2), "上传失败")
)
