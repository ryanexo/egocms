package errno

import (
    `cms/internal/erroz`
    `cms/internal/httpserver`
)

const (
    entriesDisabled = erroz.ContentType + iota
    entriesInvalidType
    entriesMissingValue
    entriesInvalidNumRange
    entriesInvalidTimeRange
    entriesInvalidLen
    entriesInvalidValue
)

var (
    ErrDisabled         = httpserver.NewError(entriesDisabled, "%s参数不可用")
    ErrInvalidValueType = httpserver.NewError(entriesInvalidType, "字段 %s 类型错误")
    ErrInvalidType      = httpserver.NewError(entriesInvalidType, "%s数据类型不合法")
    ErrMissingValue     = httpserver.NewError(entriesMissingValue, "%s必填")
    ErrInvalidNumRange  = httpserver.NewError(entriesInvalidNumRange, "%s必须在 %.2f 到 %.2f 之间")
    ErrInvalidTimeRange = httpserver.NewError(entriesInvalidTimeRange, "%s必须在 %s 到 %s 之间")
    ErrInvalidLen       = httpserver.NewError(entriesInvalidLen, "%s长度必须在 %d 到 %d 之间")
    ErrInvalidValue     = httpserver.NewError(entriesInvalidValue, "%s值不合法")
)
