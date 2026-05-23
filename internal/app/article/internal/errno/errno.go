package errno

import (
    `cms/internal/erroz`
    `cms/internal/httpserver`
)

const (
    publishStatusNotAllowed = erroz.Article + iota
    submitStatusNotAllowed
    offlineStatusNotAllowed
    rejectStatusNotAllowed
    republishStatusNotAllowed
    alreadyPublished
    alreadyRepublish
    alreadySubmitted
    alreadyOffline
    alreadyReject

)

var (
    ErrPublishStatusNotAllowed   = httpserver.NewError(publishStatusNotAllowed, "内容非待审状态，无法发布")
    ErrSubmitStatusNotAllowed    = httpserver.NewError(submitStatusNotAllowed, "内容非草稿状态，无法提交")
    ErrOfflineStatusNotAllowed   = httpserver.NewError(offlineStatusNotAllowed, "内容非发布状态，无法下线")
    ErrRejectStatusNotAllowed    = httpserver.NewError(rejectStatusNotAllowed, "内容非待审状态，无法发布")
    ErrRepublishStatusNotAllowed = httpserver.NewError(republishStatusNotAllowed, "内容非发布状态，无法提交")
    ErrAlreadyPublished          = httpserver.NewError(alreadyPublished, "文章已发布，请勿重复发布")
    ErrAlreadyRepublish          = httpserver.NewError(alreadyRepublish, "文章已提交，请勿重复提交")
    ErrAlreadySubmitted          = httpserver.NewError(alreadySubmitted, "文章已提交，请勿重复提交")
    ErrAlreadyOffline            = httpserver.NewError(alreadyOffline, "文章已下线，请勿重复下线")
    ErrAlreadyReject             = httpserver.NewError(alreadyReject, "文章已拒审，请勿重复拒审")
)
