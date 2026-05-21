package errno

import (
    `cms/internal/httpserver`
)

var (
    ErrPublishStatusNotAllowed   = httpserver.NewError(1001, "内容非待审状态，无法发布")
    ErrSubmitStatusNotAllowed    = httpserver.NewError(1002, "内容非草稿状态，无法提交")
    ErrOfflineStatusNotAllowed   = httpserver.NewError(1003, "内容非发布状态，无法下线")
    ErrRejectStatusNotAllowed    = httpserver.NewError(1004, "内容非待审状态，无法发布")
    ErrRepublishStatusNotAllowed = httpserver.NewError(1005, "内容非发布状态，无法提交")
    ErrAlreadyPublished          = httpserver.NewError(1006, "文章已发布，请勿重复发布")
    ErrAlreadyRepublish          = httpserver.NewError(1007, "文章已提交，请勿重复提交")
    ErrAlreadySubmitted          = httpserver.NewError(1008, "文章已提交，请勿重复提交")
    ErrAlreadyOffline            = httpserver.NewError(1009, "文章已下线，请勿重复下线")
    ErrAlreadyReject             = httpserver.NewError(1010, "文章已拒审，请勿重复拒审")
    
    ErrModelDataDisabled         = httpserver.NewError(1101, "%s参数不可用")
    ErrModelDataInvalidType      = httpserver.NewError(1102, "%s数据类型不合法")
    ErrModelDataMissingValue     = httpserver.NewError(1103, "%s必填")
    ErrModelDataInvalidNumRange  = httpserver.NewError(1104, "%s必须在 %.2f 到 %.2f 之间")
    ErrModelDataInvalidTimeRange = httpserver.NewError(1105, "%s必须在 %s 到 %s 之间")
    ErrModelDataInvalidLen       = httpserver.NewError(1106, "%s长度必须在 %d 到 %d 之间")
    ErrModelDataInvalidValue     = httpserver.NewError(1107, "%s值不合法")
)
