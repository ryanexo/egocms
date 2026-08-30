package errno

import (
    `cms/internal/public/erroz`
)

var (
    ErrDisabled         = erroz.NewError("CONTENT_TYPE_001", "%s参数不可用")
    ErrInvalidValueType = erroz.NewError("CONTENT_TYPE_002", "字段 %s 类型错误")
    ErrInvalidType      = erroz.NewError("CONTENT_TYPE_003", "%s数据类型不合法")
    ErrMissingValue     = erroz.NewError("CONTENT_TYPE_004", "%s必填")
    ErrInvalidNumRange  = erroz.NewError("CONTENT_TYPE_005", "%s必须在 %.2f 到 %.2f 之间")
    ErrInvalidTimeRange = erroz.NewError("CONTENT_TYPE_006", "%s必须在 %s 到 %s 之间")
    ErrInvalidLen       = erroz.NewError("CONTENT_TYPE_007", "%s长度必须在 %d 到 %d 之间")
    ErrInvalidValue     = erroz.NewError("CONTENT_TYPE_008", "%s值不合法")
)
