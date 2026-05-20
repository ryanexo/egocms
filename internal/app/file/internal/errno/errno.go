package errno

import (
    `cms/internal/erroz`
    erroz2 `cms/internal/httpserver`
)

var (
    FileDriverConfigNotExists  = erroz2.New(erroz.Code("FILE", errtype.NotFound, 0), "未配置文件驱动")
    FileDeleteFailedInCreating = erroz2.New(erroz.Code("FILE", errtype.Unknown, 1), "上传失败")
    FilePreCreateFileFailed    = erroz2.New(erroz.Code("FILE", errtype.Unknown, 2), "上传失败")
)
