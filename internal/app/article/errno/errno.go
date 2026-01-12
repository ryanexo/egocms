package errno

import (
    `dpcms/internal/erroz`
    `dpcms/internal/erroz/type`
)

var (
    ArticlePublishStatusNotAllowed   = erroz.New(erroz.Code("ARTICLE", errtype.InvalidState, 0), "内容非待审状态，无法发布")
    ArticleSubmitStatusNotAllowed    = erroz.New(erroz.Code("ARTICLE", errtype.InvalidState, 1), "内容非草稿状态，无法提交")
    ArticleOfflineStatusNotAllowed   = erroz.New(erroz.Code("ARTICLE", errtype.InvalidState, 2), "内容非发布状态，无法下线")
    ArticleRejectStatusNotAllowed    = erroz.New(erroz.Code("ARTICLE", errtype.InvalidState, 3), "内容非待审状态，无法发布")
    ArticleRepublishStatusNotAllowed = erroz.New(erroz.Code("ARTICLE", errtype.InvalidState, 4), "内容非发布状态，无法提交")
    ArticleAlreadyPublished          = erroz.New(erroz.Code("ARTICLE", errtype.Conflict, 0), "文章已发布，请勿重复发布")
    ArticleAlreadyRepublish          = erroz.New(erroz.Code("ARTICLE", errtype.Conflict, 1), "文章已提交，请勿重复提交")
    ArticleAlreadySubmitted          = erroz.New(erroz.Code("ARTICLE", errtype.Conflict, 2), "文章已提交，请勿重复提交")
    ArticleAlreadyOffline            = erroz.New(erroz.Code("ARTICLE", errtype.Conflict, 3), "文章已下线，请勿重复下线")
    ArticleAlreadyReject             = erroz.New(erroz.Code("ARTICLE", errtype.Conflict, 4), "文章已拒审，请勿重复拒审")
    
    ArticleModelDataDisabled         = erroz.New(erroz.Code("ARTICLE_MODEL", errtype.Parameter, 0), "%s参数不可用")
    ArticleModelDataInvalidType      = erroz.New(erroz.Code("ARTICLE_MODEL", errtype.Parameter, 1), "%s数据类型不合法")
    ArticleModelDataMissingValue     = erroz.New(erroz.Code("ARTICLE_MODEL", errtype.Parameter, 2), "%s必填")
    ArticleModelDataInvalidNumRange  = erroz.New(erroz.Code("ARTICLE_MODEL", errtype.Parameter, 3), "%s必须在 %.2f 到 %.2f 之间")
    ArticleModelDataInvalidTimeRange = erroz.New(erroz.Code("ARTICLE_MODEL", errtype.Parameter, 3), "%s必须在 %s 到 %s 之间")
    ArticleModelDataInvalidLen       = erroz.New(erroz.Code("ARTICLE_MODEL", errtype.Parameter, 4), "%s长度必须在 %d 到 %d 之间")
    ArticleModelDataInvalidValue     = erroz.New(erroz.Code("ARTICLE_MODEL", errtype.Parameter, 5), "%s值不合法")
)
