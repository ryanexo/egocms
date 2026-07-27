package errno

import (
    `cms/internal/httpx`
)

var (
    ErrDisabled         = httpx.NewError(20001, "%s参数不可用")
    ErrInvalidValueType = httpx.NewError(20002, "字段 %s 类型错误")
    ErrInvalidType      = httpx.NewError(20003, "%s数据类型不合法")
    ErrMissingValue     = httpx.NewError(20004, "%s必填")
    ErrInvalidNumRange  = httpx.NewError(20005, "%s必须在 %.2f 到 %.2f 之间")
    ErrInvalidTimeRange = httpx.NewError(20006, "%s必须在 %s 到 %s 之间")
    ErrInvalidLen       = httpx.NewError(20007, "%s长度必须在 %d 到 %d 之间")
    ErrInvalidValue     = httpx.NewError(20008, "%s值不合法")
)
